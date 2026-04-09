package cmd

import (
	"print-agent/internal/adapter/api/printer"
	"print-agent/internal/adapter/api/ws/handler"
	"print-agent/internal/adapter/database"
	"print-agent/internal/adapter/tui"
	"print-agent/internal/core/invoice"
	printer2 "print-agent/internal/core/printer"
	"print-agent/internal/core/setting"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var Singleton *App

type App struct {
	Tui        Tui
	Handlers   Handlers
	Usecase    Usecase
	Repository Repository
	Route      Route
}

func (app *App) Bootstrap(fiber *fiber.App, db *gorm.DB) {
	printers := printer.NewPrinter()

	app.Repository = *app.NewRepositories()
	app.Usecase = *app.NewUsecases(db, printers)
	app.Handlers = *app.NewHandlers(app.Usecase)
	app.Tui = *app.NewTui(printers)
	//app.Route = *app.NewRoutes(fiber)

	Singleton = app
}

type Tui struct {
	DiscoveryForm  tui.DiscoveryForm
	CompanyForm    tui.CompanyForm
	ConnectionForm tui.ConnectionForm
	DeviceForm     tui.DeviceForm
}

func (app *App) NewTui(printers printer2.Printer) *Tui {
	return &Tui{
		DiscoveryForm:  *tui.NewDiscoveryForm(printers, app.Usecase.SettingUsecase),
		CompanyForm:    *tui.NewCompanyForm(app.Usecase.SettingUsecase),
		ConnectionForm: *tui.NewConnectionForm(app.Usecase.SettingUsecase),
		DeviceForm:     *tui.NewDeviceForm(app.Usecase.SettingUsecase),
	}
}

type Handlers struct {
	ReceiverHandler handler.ReceiverHandler
}

func (app *App) NewHandlers(usecase Usecase) *Handlers {

	return &Handlers{
		ReceiverHandler: *handler.NewReceiverHandler(usecase.InvoiceUsecase, usecase.SettingUsecase),
	}
}

type Usecase struct {
	SettingUsecase setting.Usecase
	InvoiceUsecase invoice.Usecase
}

func (app *App) NewUsecases(db *gorm.DB, printers printer2.Printer) *Usecase {
	settingDependency := setting.UsecaseDependency{
		DB:                db,
		SettingRepository: app.Repository.SettingRepository,
	}

	settingUsecase := setting.NewUsecase(settingDependency)

	invoiceDependency := invoice.UsecaseDependency{
		DB:             db,
		SettingUsecase: settingUsecase,
		Printer:        printers,
	}
	return &Usecase{
		SettingUsecase: settingUsecase,
		InvoiceUsecase: invoice.NewUsecase(invoiceDependency),
	}
}

type Repository struct {
	SettingRepository setting.Repository
}

func (app *App) NewRepositories() *Repository {
	return &Repository{
		SettingRepository: database.NewSettingRepository(),
	}
}

type Route struct {
	ReceiverHandler handler.ReceiverHandler
}

//func (app *App) NewRoutes(fiber *fiber.App) *Route {
//	userRoute := *route.NewUserRoutes(&app.Handlers.UserHandler)
//	authRoute := *route.NewAuthRoutes(&app.Handlers.AuthHandler)
//
//	router := fiber.Group("/api/v1")
//	userRoute.InstallRoutes(router)
//	authRoute.InstallRoutes(router)
//
//	return &Route{
//		UserRoute: userRoute,
//		AuthRoute: authRoute,
//	}
//}
