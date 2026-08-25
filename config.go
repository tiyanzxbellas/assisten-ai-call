package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GroqAPI      string
	GeminiAPI    string
	GeminiModel  string
	SystemPrompt string
	PromptLang   string
	TTSCharacter string
	TTSLang      string
	TTSSpeed     float64
	Owners       []string
}

var AppConfig *Config

const DefaultSystemPrompt = `Kamu adalah asisten suara AI pintar yang sedang berbicara lewat panggilan telepon WhatsApp.
Aturan penting:
1. Jawablah dengan singkat, ramah, jujur, dan alami seperti percakapan telepon manusia nyata (maksimal 2-3 kalimat per respon).
2. JANGAN gunakan format markdown seperti bintang (*), pagar (#), bullet point (-), atau simbol tabel, karena responmu akan langsung diubah menjadi suara (Text-to-Speech).
3. Gunakan Bahasa Indonesia yang santai dan sopan.`

func LoadConfig() *Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	groqAPI := os.Getenv("GROQ_API")
	if groqAPI == "" {
		groqAPI = os.Getenv("GROQ_API_KEY")
	}

	geminiAPI := os.Getenv("GEMINI_API")
	if geminiAPI == "" {
		geminiAPI = os.Getenv("GEMINI_API_KEY")
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		geminiModel = "gemini-3.6-flash"
	}

	sysPrompt := os.Getenv("SYSTEM_PROMPT")
	if sysPrompt == "" {
		sysPrompt = DefaultSystemPrompt
	}

	promptLang := os.Getenv("PROMPT_LANG")
	if promptLang == "" {
		promptLang = "id"
	}

	ttsChar := os.Getenv("TTS_CHARACTER")
	if ttsChar == "" {
		ttsChar = "特别周 Special Week (Umamusume Pretty Derby)"
	}

	ttsLang := os.Getenv("TTS_LANG")
	if ttsLang == "" {
		ttsLang = "Mix"
	}

	ttsSpeedStr := os.Getenv("TTS_SPEED")
	ttsSpeed := 1.0
	if ttsSpeedStr != "" {
		if val, err := strconv.ParseFloat(ttsSpeedStr, 64); err == nil && val > 0 {
			ttsSpeed = val
		}
	}

	ownerStr := os.Getenv("OWNER")
	var owners []string
	if ownerStr != "" {
		parts := strings.Split(ownerStr, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				owners = append(owners, p)
			}
		}
	}

	AppConfig = &Config{
		GroqAPI:      groqAPI,
		GeminiAPI:    geminiAPI,
		GeminiModel:  geminiModel,
		SystemPrompt: sysPrompt,
		PromptLang:   promptLang,
		TTSCharacter: ttsChar,
		TTSLang:      ttsLang,
		TTSSpeed:     ttsSpeed,
		Owners:       owners,
	}

	if AppConfig.GroqAPI == "" {
		log.Println("[WARNING] GROQ_API key is missing in .env! STT functionality might fail.")
	}
	if AppConfig.GeminiAPI == "" {
		log.Println("[WARNING] GEMINI_API key is missing in .env! AI chat functionality might fail.")
	}

	tempDir := filepath.Join(".", "temp")
	_ = os.MkdirAll(tempDir, 0755)

	return AppConfig
}

func IsOwner(sender string) bool {
	if len(AppConfig.Owners) == 0 {
		return true
	}
	cleanedSender := strings.TrimPrefix(sender, "+")
	cleanedSender = strings.Split(cleanedSender, "@")[0]
	cleanedSender = strings.Split(cleanedSender, ".")[0]
	cleanedSender = strings.TrimSpace(cleanedSender)

	for _, owner := range AppConfig.Owners {
		cleanedOwner := strings.TrimPrefix(owner, "+")
		cleanedOwner = strings.Split(cleanedOwner, "@")[0]
		cleanedOwner = strings.TrimSpace(cleanedOwner)
		if cleanedOwner == "" {
			continue
		}
		if cleanedSender == cleanedOwner || strings.HasSuffix(cleanedSender, cleanedOwner) || strings.HasPrefix(cleanedSender, cleanedOwner) {
			return true
		}
	}
	return false
}
