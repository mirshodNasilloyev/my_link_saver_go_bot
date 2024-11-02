package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/config"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/repository"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/repository/boltdb"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/server"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("Initializing config error", err)
	}
	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("Incorrect token: ", err)
	}
	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PockectConsumerKey)
	if err != nil {
		log.Fatal("Incorrect pocket client: ", err)
	}
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Incorrect db: ", err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)
	authorizedServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizedServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}
		return nil
	})
	return db, nil
}
