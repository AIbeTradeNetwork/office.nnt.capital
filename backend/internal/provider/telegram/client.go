package telegram

import (
	"github.com/mr-linch/go-tg"
	"server/internal/config"
)

func NewClient() *tg.Client {
	cfg := config.Get()

	return tg.New(cfg.TgBotApiKey)
}
