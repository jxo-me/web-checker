package main

import (
	"github.com/jxo-me/web-checker/config"
	"github.com/jxo-me/web-checker/core"
	"github.com/jxo-me/web-checker/hook"
	"log"
)

func main() {
	cfg, err := config.FromEnv("CONFIG")
	if err != nil {
		log.Fatal(err)
	}

	telegramHook := &hook.TelegramHook{Token: cfg.Telegram.Token, ChatId: cfg.Telegram.ChatId, TimeOut: cfg.Checker.Timeout}
	err = telegramHook.Notify(cfg.Checker.Websites)
	if err != nil {
		log.Fatal(err)
	}
	checker := &core.Checker{
		Config: cfg.Checker,
		Processors: core.Processors{
			telegramHook.Process,
		},
	}

	log.Println("run the website checker")
	checker.Run()

}
