package tui

import (
	"context"
	"print-agent/internal/core/setting"

	"github.com/charmbracelet/huh"
)

type CompanyForm struct {
	SettingUsecase setting.Usecase
}

func NewCompanyForm(setting setting.Usecase) *CompanyForm {
	return &CompanyForm{
		SettingUsecase: setting,
	}
}

func (tui *CompanyForm) Show() {
	var storeName string
	var address string
	err :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Set your store name").
					Value(&storeName),
				huh.NewInput().
					Title("Set your store address").
					Value(&address))).
			Run()

	if err != nil {
		return
	}

	request := setting.Content{
		StoreName: storeName,
		Address:   address,
	}
	_, err = tui.SettingUsecase.Save(context.Background(), request)
	if err != nil {
		return
	}
}
