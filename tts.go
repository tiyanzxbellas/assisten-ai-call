package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const HF_BASE_URL = "https://plachta-vits-umamusume-voice-synthesizer.hf.space"

type GradioJoinRequest struct {
	Data        []interface{} `json:"data"`
	EventData   interface{}   `json:"event_data"`
	FnIndex     int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
	TriggerID   int           `json:"trigger_id"`
}

type GradioSSEMessage struct {
	Msg    string `json:"msg"`
	Output struct {
		Data []interface{} `json:"data"`
	} `json:"output"`
}

func randomSessionHash() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateTTS converts text into an audio file (WAV/MP3) path.
// Tries Anime TTS (HF space) first, with fallback to Google Translate TTS.
func GenerateTTS(text string) (string, error) {
	sessionHash := randomSessionHash()
	char := AppConfig.TTSCharacter
	if char == "" {
		char = "特别周 Special Week (Umamusume Pretty Derby)"
	}
	lang := AppConfig.TTSLang
	if lang == "" {
		lang = "Mix"
	}

	tempOutPath := filepath.Join(".", "temp", fmt.Sprintf("tts_%s.wav", sessionHash))

	// Attempt Anime TTS via HuggingFace Gradio Queue API
	audioPath, err := generateAnimeTTS(text, char, lang, sessionHash, tempOutPath)
	if err == nil && audioPath != "" {
		return audioPath, nil
	}

	fmt.Printf("[TTS] Anime TTS failed/timed out (%v). Using Google Translate TTS fallback...\n", err)
	return generateGoogleTTS(text, filepath.Join(".", "temp", fmt.Sprintf("tts_fallback_%s.mp3", sessionHash)))
}

func generateAnimeTTS(text, char, lang, sessionHash, outputPath string) (string, error) {
	client := &http.Client{Timeout: 45 * time.Second}

	joinReq := GradioJoinRequest{
		Data: []interface{}{
			text,
			char,
			lang,
			AppConfig.TTSSpeed, // speed
			false,              // noise
		},
		EventData:   nil,
		FnIndex:     2,
		SessionHash: sessionHash,
		TriggerID:   24,
	}

	jsonBytes, err := json.Marshal(joinReq)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(HF_BASE_URL+"/gradio_api/queue/join", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("queue join failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("queue join returned status %d: %s", resp.StatusCode, string(b))
	}

	pollURL := fmt.Sprintf("%s/gradio_api/queue/data?session_hash=%s", HF_BASE_URL, sessionHash)
	deadline := time.Now().Add(35 * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(1500 * time.Millisecond)

		req, _ := http.NewRequest("GET", pollURL, nil)
		pollResp, err := client.Do(req)
		if err != nil {
			continue
		}

		bodyBytes, err := io.ReadAll(pollResp.Body)
		pollResp.Body.Close()
		if err != nil {
			continue
		}

		lines := strings.Split(string(bodyBytes), "\n")
		for _, line := range lines {
			if !strings.HasPrefix(line, "data:") {
				continue
			}

			jsonData := strings.TrimPrefix(line, "data:")
			jsonData = strings.TrimSpace(jsonData)

			var msg GradioSSEMessage
			if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
				continue
			}

			if msg.Msg == "process_completed" {
				for _, item := range msg.Output.Data {
					var fileURL string

					switch v := item.(type) {
					case map[string]interface{}:
						if urlVal, ok := v["url"].(string); ok && urlVal != "" {
							fileURL = urlVal
						} else if nameVal, ok := v["name"].(string); ok && nameVal != "" {
							fileURL = HF_BASE_URL + "/file=" + nameVal
						}
					case string:
						if strings.HasPrefix(v, "http") || strings.HasPrefix(v, "/file=") {
							fileURL = v
						}
					}

					if fileURL != "" {
						if !strings.HasPrefix(fileURL, "http") {
							fileURL = HF_BASE_URL + fileURL
						}
						err := downloadFile(fileURL, outputPath)
						if err == nil {
							return outputPath, nil
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("Anime TTS processing timeout")
}

func generateGoogleTTS(text, outputPath string) (string, error) {
	encodedText := url.QueryEscape(text)
	ttsURL := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=id&client=tw-ob&q=%s", encodedText)

	err := downloadFile(ttsURL, outputPath)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

func downloadFile(fileURL, outputPath string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
