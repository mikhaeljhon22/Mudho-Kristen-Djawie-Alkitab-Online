//go:build wireinject
// +build wireinject

package main

import (
	"alkitab/configs"
	"alkitab/controllers"
	"alkitab/entitys"
	"alkitab/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func ProvideSQLConnect() *gorm.DB {
	db := configs.ConnectPostgre()
	db.AutoMigrate(&entitys.UsersLetstalk{})
	return db
}

func ProvideUserService(db *gorm.DB) *services.UserService {
	return services.NewUserService(db)
}

func ProvideMailService() *services.MailService {
	return services.NewMailService()
}

func ProvideEmailWorker(mailService *services.MailService) chan services.EmailJob {
	return services.StartEmailWorkerPool(mailService, 5)
}

func ProvideUserController(userService *services.UserService, mailService *services.MailService, emailJob chan services.EmailJob) *controllers.UserController {
	ctrl := controllers.NewUserController(userService, mailService)
	ctrl.EmailJobs = emailJob
	return ctrl
}

type Server struct {
	UserController *controllers.UserController
}

func InitializedServer() *Server {
	wire.Build(
		ProvideSQLConnect,
		ProvideUserService,
		ProvideMailService,
		ProvideEmailWorker,
		ProvideUserController,
		wire.Struct(new(Server), "UserController"),
	)
	return nil
}
