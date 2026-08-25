package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/purpshell/meowcaller"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

func main() {
	log.Println("==================================================")
	log.Println("     WhatsApp AI Call Assistant Bot (Golang)      ")
	log.Println("==================================================")

	// Load configuration from .env
	cfg := LoadConfig()
	log.Printf("[Config] GROQ_API status: %t", cfg.GroqAPI != "")
	log.Printf("[Config] GEMINI_API status: %t", cfg.GeminiAPI != "")
	log.Printf("[Config] GEMINI_MODEL: %s", cfg.GeminiModel)
	log.Printf("[Config] TTS Character: %s", cfg.TTSCharacter)
	log.Printf("[Config] TTS Speed: %.1f", cfg.TTSSpeed)
	if len(cfg.Owners) > 0 {
		log.Printf("[Config] Owners configured: %v", cfg.Owners)
	} else {
		log.Println("[Config] Owners: Public (Anyone can call)")
	}

	// Setup SQLite container store for WhatsApp session
	dbLog := waLog.Stdout("Database", "WARN", true)
	container, err := sqlstore.New(context.Background(), "sqlite", "file:ai_call_session.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get device store: %v", err)
	}

	// Create whatsmeow client
	clientLog := waLog.Stdout("Client", "INFO", true)
	waClient := whatsmeow.NewClient(deviceStore, clientLog)

	// Initialize meowcaller client
	callerClient := meowcaller.NewClient(waClient)

	// Setup Incoming Call Handler
	callerClient.OnIncomingCall(func(call *meowcaller.Call) {
		peer := call.Peer().User
		log.Printf("[Call] Incoming call from %s (Call ID: %s)", call.Peer().String(), call.ID())

		if !IsOwner(peer) {
			log.Printf("[Call] Rejecting incoming call from non-owner: %s", peer)
			_ = call.Reject()
			return
		}

		// Answer incoming call
		if err := call.Answer(); err != nil {
			log.Printf("[Call] Failed to answer call: %v", err)
			return
		}

		session := NewAICallSession(call)
		call.OnEnd(func(reason string) {
			log.Printf("[Call] Call ended (Reason: %s)", reason)
			session.Stop()
		})

		// Start AI Voice Loop in background
		go session.StartVoiceLoop()
	})

	// Setup Message Handler for commands (e.g., !aicall 628xxx)
	waClient.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			handleMessage(waClient, callerClient, v)
		}
	})

	// QR Code / Session Login
	if waClient.Store.ID == nil {
		qrChan, _ := waClient.GetQRChannel(context.Background())
		err = waClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("\nScan QR Code dibawah ini menggunakan WhatsApp di HP Anda:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Printf("QR Event: %s\n", evt.Event)
			}
		}
	} else {
		err = waClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		log.Println("Berhasil terhubung ke WhatsApp!")
	}

	// Keep program running until interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down AI Call Assistant...")
	waClient.Disconnect()
}

func handleMessage(client *whatsmeow.Client, callerClient *meowcaller.Client, msg *events.Message) {
	if msg.Info.IsFromMe {
		return
	}

	body := getMessageBody(msg.Message)
	if body == "" {
		return
	}

	fields := strings.Fields(body)
	if len(fields) == 0 {
		return
	}
	cmd := strings.ToLower(fields[0])

	if cmd == "!aicall" || cmd == ".aicall" || cmd == "!call" || cmd == ".call" {
		senderUser := msg.Info.Sender.User
		if !IsOwner(senderUser) {
			log.Printf("[Command] Denied call command from non-owner: %s", senderUser)
			_, _ = client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
				Conversation: buildStringPointer("Maaf, fitur ini hanya dapat digunakan oleh Owner bot."),
			})
			return
		}

		if len(fields) < 2 {
			_, _ = client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
				Conversation: buildStringPointer("Gunakan format: !aicall 628xxxxxxxx"),
			})
			return
		}

		targetNum := fields[1]
		targetNum = strings.TrimPrefix(targetNum, "+")
		targetNum = strings.ReplaceAll(targetNum, "-", "")
		targetNum = strings.ReplaceAll(targetNum, " ", "")

		if targetNum == "" {
			_, _ = client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
				Conversation: buildStringPointer("Gunakan format: !aicall 628xxxxxxxx"),
			})
			return
		}

		targetJID := types.NewJID(targetNum, types.DefaultUserServer)
		_, _ = client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
			Conversation: buildStringPointer(fmt.Sprintf("Memulai panggilan AI ke nomor %s...", targetNum)),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		log.Printf("[Command] Dialing %s for AI call...", targetJID.String())
		call, err := callerClient.Call(ctx, targetJID.String())
		if err != nil {
			log.Printf("[Command] Failed to place call: %v", err)
			_, _ = client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
				Conversation: buildStringPointer(fmt.Sprintf("Gagal menelpon: %v", err)),
			})
			return
		}

		session := NewAICallSession(call)
		call.OnEnd(func(reason string) {
			log.Printf("[Call] Call with %s ended: %s", targetJID.String(), reason)
			session.Stop()
		})

		call.OnReady(func() {
			log.Printf("[Call] Call connected with %s! Starting AI voice loop...", targetJID.String())
			go session.StartVoiceLoop()
		})
	}
}

func getMessageBody(m *waProto.Message) string {
	if m == nil {
		return ""
	}
	if m.GetEphemeralMessage() != nil {
		return getMessageBody(m.GetEphemeralMessage().GetMessage())
	}
	if m.GetViewOnceMessage() != nil {
		return getMessageBody(m.GetViewOnceMessage().GetMessage())
	}
	if m.GetViewOnceMessageV2() != nil {
		return getMessageBody(m.GetViewOnceMessageV2().GetMessage())
	}
	if m.GetDocumentWithCaptionMessage() != nil {
		return getMessageBody(m.GetDocumentWithCaptionMessage().GetMessage())
	}
	if m.GetConversation() != "" {
		return m.GetConversation()
	}
	if m.GetExtendedTextMessage() != nil {
		return m.GetExtendedTextMessage().GetText()
	}
	return ""
}

func buildStringPointer(s string) *string {
	return &s
}
