package setting

import (
	"regexp"
	"strings"
)

type Content struct {
	DeviceID       string
	ChoosenPrinter string
	WebsocketURL   string
	WebsocketKey   string
	ChannelName    string
	StoreName      string
	Address        string
}

func (dto *Content) ToEntity() *Setting {
	return &Setting{
		DeviceID:       toKebabCase(dto.DeviceID),
		ChoosenPrinter: dto.ChoosenPrinter,
		WebsocketURL:   dto.WebsocketURL,
		WebsocketKey:   dto.WebsocketKey,
		ChannelName:    dto.ChannelName,
		StoreName:      dto.StoreName,
		Address:        dto.Address,
	}
}

func ToContent(entity *Setting) *Content {
	return &Content{
		DeviceID:       entity.DeviceID,
		ChoosenPrinter: entity.ChoosenPrinter,
		WebsocketURL:   entity.WebsocketURL,
		WebsocketKey:   entity.WebsocketKey,
		ChannelName:    entity.ChannelName,
		StoreName:      entity.StoreName,
		Address:        entity.Address,
	}
}

func toKebabCase(s string) string {
	// Pisahin camelCase / PascalCase jadi spasi
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re.ReplaceAllString(s, "${1} ${2}")

	// Replace non-alphanumeric jadi spasi
	re = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	s = re.ReplaceAllString(s, " ")

	// Trim, to lower, replace spasi jadi dash
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", "-"))
}
