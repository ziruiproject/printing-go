package tui

import (
	"context"
	"print-agent/internal/core/setting"

	"github.com/charmbracelet/huh"
)

type DeviceForm struct {
	SettingUsecase setting.Usecase
}

func NewDeviceForm(setting setting.Usecase) *DeviceForm {
	return &DeviceForm{
		SettingUsecase: setting,
	}
}

func (tui *DeviceForm) Show() {
	var deviceId string
	err := huh.NewInput().
		Title("Set your device id").
		Value(&deviceId).
		Description("Note: device id will be formatted to kebab-case").
		Run()

	if err != nil {
		return
	}

	request := setting.Content{
		DeviceID: deviceId,
	}
	_, err = tui.SettingUsecase.Save(context.Background(), request)
	if err != nil {
		return
	}
}
