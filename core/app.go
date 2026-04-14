package core

import (
	"github.com/nanoteck137/refinery/config"
	"github.com/nanoteck137/refinery/database"
	"github.com/nanoteck137/refinery/service"
	"github.com/nanoteck137/refinery/tools/broker"
	"github.com/nanoteck137/refinery/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	DB() *database.Database
	Config() *config.Config

	DataDir() types.DataDir

	NotificationService() *service.NotificationService
	TaskService() *service.TaskService
	AuthService() *service.AuthService
	UserService() *service.UserService
	ImageService() *service.ImageService

	Broker() *broker.Broker

	Bootstrap() error
	Shutdown() error
}
