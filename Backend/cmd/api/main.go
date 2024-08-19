package main

import (
	"fmt"
	"log"

	"github.com/Powerdfi-com/Backend/cmd/api/routes"
	"github.com/Powerdfi-com/Backend/config"
	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	"github.com/Powerdfi-com/Backend/external/shufti"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/cron/minter"
	"github.com/Powerdfi-com/Backend/internal/cron/orderbook"
	"github.com/Powerdfi-com/Backend/internal/implementation/postgres"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	cfg, err := config.LoadConfig()
	failOnError(err, "error parsing config")

	db, err := openDB(cfg.DbUri)
	failOnError(err, "error parsing config")
	defer db.Close()
	log.Println("database connection established")

	hederaClient, err := hederaUtils.CreateClient(cfg)
	failOnError(err, "error connecting to hedera")
	log.Println("hedera connection established")

	repos := postgres.NewRepositories(db)

	app := internal.Application{
		Config:       cfg,
		Repositories: repos,
		HederaClient: hederaClient,
		ShuftiClient: shufti.NewShuftiClient(
			cfg.Shufti.BaseUrl,
			cfg.Shufti.ClientId,
			cfg.Shufti.SecretKey,
			cfg.Shufti.JourneyId,
		),
	}

	engine := routes.Engine(app)

	m := minter.NewMinter(app)
	go m.WatchDbForApprovedAssets()

	o := orderbook.NewOrderBook(app)
	go o.WatchDbOrderMatches()
	// start server
	if err := engine.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Panicf("%s: %s", "server stopped", err)
	}

	log.Println("server stopped")
}
