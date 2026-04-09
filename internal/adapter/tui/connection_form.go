package tui

import (
	"context"
	"print-agent/internal/core/setting"

	"github.com/charmbracelet/huh"
)

type ConnectionForm struct {
	SettingUsecase setting.Usecase
}

func NewConnectionForm(setting setting.Usecase) *ConnectionForm {
	return &ConnectionForm{
		SettingUsecase: setting,
	}
}

func (tui *ConnectionForm) Show() {
	var websocketUrl string
	var websocketKey string
	var channelId string
	err :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Set websocket url").
					Value(&websocketUrl),
				huh.NewInput().
					Title("Set websocket key").
					Value(&websocketKey),
				huh.NewInput().
					Title("Set channel id").
					Value(&channelId))).
			Run()

	if err != nil {
		return
	}

	request := setting.Content{
		WebsocketURL: websocketUrl,
		WebsocketKey: websocketKey,
		ChannelName:  channelId,
	}
	_, err = tui.SettingUsecase.Save(context.Background(), request)
	if err != nil {
		return
	}
}
