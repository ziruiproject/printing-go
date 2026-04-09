package setting

import "gorm.io/gorm"

type Setting struct {
	DeviceID       string
	ChoosenPrinter string
	WebsocketURL   string
	WebsocketKey   string
	ChannelName    string
	StoreName      string
	Address        string
	gorm.Model
}
