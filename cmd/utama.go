package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Lebar ideal untuk printer 80mm
	// Sesuaikan nilai ini jika cetakan masih tidak rapi
	printerWidth = 48
)

// Definisi struct untuk data yang diterima dari Laravel
// Pastikan nama field sesuai dengan JSON dari Laravel

type BroadcastEnvelope struct {
	Event   string `json:"event"`
	Data    string `json:"data"`
	Channel string `json:"channel"`
}

type PrintData struct {
	CustomerName string  `json:"customer_name"`
	InvoiceID    string  `json:"invoice_id"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
	Discount     float64 `json:"discount"`
	GrandTotal   float64 `json:"grand_total"`
}

type Item struct {
	Description string  `json:"description"`
	Qty         int     `json:"qty"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

// padRight mengisi spasi di sebelah kanan string hingga mencapai lebar tertentu
func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// padLeft mengisi spasi di sebelah kiri string hingga mencapai lebar tertentu
func padLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

// centerText menengahkan string dengan menambahkan spasi di kedua sisi
func centerText(s string, length int) string {
	if len(s) >= length {
		return s
	}
	padding := (length - len(s)) / 2
	return strings.Repeat(" ", padding) + s + strings.Repeat(" ", length-len(s)-padding)
}

func main() {
	printerPath := "/dev/usb/lp0"
	deviceID := "printer-komputer-1"
	wsURL := "ws://localhost:8080/app/" + "luhtvhd4jpjnkft0tlmm" + "?device_id=" + deviceID

	// Logika koneksi WebSocket
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatalf("URL parsing error: %v", err)
	}

	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("WebSocket connection error: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to WebSocket server. Waiting for print jobs...")

	// 1. Buat pesan 'pusher:subscribe'
	subscribeMsg := map[string]interface{}{
		"event": "pusher:subscribe",
		"data": map[string]string{
			// Channel harus sesuai dengan yang didefinisikan di Laravel
			"channel": "print-channel." + deviceID,
		},
	}

	// 2. Kirim pesan dalam format JSON
	subscribePayload, err := json.Marshal(subscribeMsg)
	if err != nil {
		log.Fatalf("Error marshalling subscribe message: %v", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, subscribePayload); err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	log.Printf("Successfully subscribed to channel: %s", "print-channel."+deviceID)

	// Loop tak terbatas untuk mendengarkan pesan dari WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("Received message: %s", message)

		var envelope BroadcastEnvelope
		err = json.Unmarshal(message, &envelope)
		if err != nil {
			log.Println("Error parsing envelope JSON:", err)
			continue
		}

		var data PrintData
		err = json.Unmarshal([]byte(envelope.Data), &data)
		if err != nil {
			log.Println("Error parsing inner data JSON:", err)
			continue
		}

		fmt.Printf("%+v\n", data)

		if data.InvoiceID == "" {
			continue
		}

		// Membangun konten cetak dari data JSON
		content := buildPrintContent(data)

		// Mencetak konten
		printToPrinter(printerPath, content)
	}
}

// Fungsi untuk membangun string konten cetak dari data JSON
func buildPrintContent(data PrintData) string {
	var sb strings.Builder

	// Header
	sb.WriteString(centerText("NOTA PEMBAYARAN", printerWidth) + "\n")
	sb.WriteString(strings.Repeat("=", printerWidth) + "\n")

	// Detail Transaksi
	sb.WriteString(padRight("Nama Pelanggan: "+data.CustomerName, printerWidth) + "\n")
	sb.WriteString(padRight("No. Faktur    : "+data.InvoiceID, printerWidth) + "\n")
	sb.WriteString(strings.Repeat("-", printerWidth) + "\n")

	// Item
	sb.WriteString(padRight("Deskripsi", 18) + "Qty  Harga    Total\n")
	sb.WriteString(strings.Repeat("-", printerWidth) + "\n")

	for _, item := range data.Items {
		// Format harga dan total ke string
		priceStr := strconv.FormatFloat(item.Price, 'f', 0, 64)
		subtotalStr := strconv.FormatFloat(item.Subtotal, 'f', 0, 64)
		qtyStr := strconv.Itoa(item.Qty)

		itemLine := padRight(item.Description, 18) + padLeft(qtyStr, 3) + " " + padLeft(priceStr, 6) + " " + padLeft(subtotalStr, 6)
		sb.WriteString(itemLine + "\n")
	}

	sb.WriteString(strings.Repeat("-", printerWidth) + "\n")

	// Total
	sb.WriteString(padRight("Total Pembelian:", 25) + padLeft("Rp"+strconv.FormatFloat(data.Total, 'f', 0, 64), 16) + "\n")
	sb.WriteString(padRight("Diskon (10%):", 25) + padLeft("Rp"+strconv.FormatFloat(data.Discount, 'f', 0, 64), 16) + "\n")
	sb.WriteString(strings.Repeat("-", printerWidth) + "\n")
	sb.WriteString(padRight("Total Bayar:", 25) + padLeft("Rp"+strconv.FormatFloat(data.GrandTotal, 'f', 0, 64), 16) + "\n")

	// Footer
	sb.WriteString(strings.Repeat("=", printerWidth) + "\n")
	sb.WriteString(centerText("Terima kasih telah berbelanja!", printerWidth) + "\n")
	sb.WriteString("\n\n\f") // Tambahkan form feed di akhir

	return sb.String()
}

// Fungsi untuk mencetak konten ke printer
func printToPrinter(printerPath, content string) {
	f, err := os.OpenFile(printerPath, os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Error opening printer device '%s': %v", printerPath, err)
		return
	}
	defer f.Close()

	bufferSize := 128
	data := []byte(content)

	for i := 0; i < len(data); i += bufferSize {
		end := i + bufferSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		_, err := f.Write(chunk)
		if err != nil {
			log.Printf("Error writing to printer: %v", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("Berhasil mencetak ke printer '%s'", printerPath)
}
