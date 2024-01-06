package main

import (
	"discount/db"
	"discount/handler"
	"discount/internal/config"
	"discount/internal/locale"
	"discount/internal/logger"
	"discount/server"
	giftService "discount/service/gift"
	discountStorage "discount/storage/discount"
	giftStorage "discount/storage/gift"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			postgresDB,
			redisDB,

			// Storages
			giftStorage.New,
			discountStorage.New,

			// services
			giftService.New,
			//TODO add discount service too

			// handlers
			handler.NewGiftHandler,
			//TODO add discount handler too

			server.NewServer,
		),
		fx.Supply(),
		fx.Invoke(
			config.Init,
			logger.SetupLogger,
			locale.Init,
			db.Migrate,
			setupServer,
			handler.SetupGiftRoutes,
			server.Run,
		),
	).Run()
}
