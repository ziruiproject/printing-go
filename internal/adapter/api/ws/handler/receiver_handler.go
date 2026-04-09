package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"print-agent/config"
	"print-agent/internal/core/invoice"
	"print-agent/internal/core/setting"
	"time"

	"github.com/gorilla/websocket"
)

type ReceiverHandler struct {
	InvoiceUsecase invoice.Usecase
	SettingUsecase setting.Usecase
}

func NewReceiverHandler(usecase invoice.Usecase, setting setting.Usecase) *ReceiverHandler {
	return &ReceiverHandler{
		InvoiceUsecase: usecase,
		SettingUsecase: setting,
	}
}

// Ini method buat connect ke server WebSocket (CLIENT)
func (handler *ReceiverHandler) StartWebSocketClient() {
	appConfig := config.AppConfig

	conf, err := handler.SettingUsecase.Find(context.Background())
	if err != nil {
		log.Println("Gagal ambil setting:", err)
		return
	}

	u := url.URL{
		Scheme: appConfig.Protocol, // atau "ws" kalau non-SSL
		Host:   conf.WebsocketURL,  // misalnya "ws.example.com:443"
		Path:   fmt.Sprintf("/app/%s", conf.WebsocketKey),
	}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	channelName := fmt.Sprintf("%s.%s", conf.ChannelName, conf.DeviceID)

	subscribeMsg := map[string]interface{}{
		"event": "pusher:subscribe",
		"data": map[string]string{
			"channel": channelName,
		},
	}

	subscribeJSON, _ := json.Marshal(subscribeMsg)
	if err := conn.WriteMessage(websocket.TextMessage, subscribeJSON); err != nil {
		log.Println("Gagal kirim subscribe:", err)
		return
	}
	log.Println("Berhasil subscribe ke channel:", channelName)

	// loop nerima pesan
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Pesan diterima: %s", msg)

		var envelope invoice.BroadcastEnvelope
		if err := json.Unmarshal(msg, &envelope); err != nil {
			log.Printf("Gagal unmarshal envelope: %v", err)
			continue
		}

		if envelope.Event == "pusher:ping" {
			pong := map[string]string{"event": "pusher:pong"}
			pongBytes, _ := json.Marshal(pong)
			conn.WriteMessage(websocket.TextMessage, pongBytes)
			log.Println("Pong terkirim ✅")
			continue
		}

		// --- data masih string JSON, jadi kita decode lagi ---
		var raw string
		if err := json.Unmarshal(envelope.Data, &raw); err != nil {
			log.Printf("Data bukan string, coba langsung object: %v", err)
			raw = string(envelope.Data) // fallback kalau ternyata bukan string
		}

		var content invoice.Content
		if err := json.Unmarshal([]byte(raw), &content); err != nil {
			log.Printf("Gagal unmarshal content: %v", err)
			continue
		}

		log.Printf("Content diterima: %+v", content)

		if content.InvoiceID == "" {
			continue
		}

		// Oper ke InvoiceUsecase
		if err := handler.InvoiceUsecase.BuildInvoice(context.Background(), content); err != nil {
			log.Println("Gagal proses invoice:", err)
		}
	}

	log.Println("Reconnecting in 5s...")
	time.Sleep(5 * time.Second)
}
