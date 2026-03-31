package main

import (
	"flag"
	"log"

	"github.com/Aidajy111/Read-adviser-bot/clients/telegramm"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegramm.NewClient(tgBotHost, mustToken())
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for acces to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}
}
