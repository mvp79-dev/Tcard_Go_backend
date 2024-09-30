package main

import (
	"t-card/bootstrap"
	"t-card/bot"
	chat "t-card/internal/bot"
	"t-card/internal/config"
)

func main() {
	cfg := config.MustLoad()


	tgBot := bot.TgBot{
		Bot: bot.InitBot(cfg.BotToken),
		BotID: cfg.BotID,
		ChannelID: cfg.ChannelID,
		AdminUserID: cfg.AdminUserID,
		TmaURL: cfg.TmaURL,
	}

	go func() {
		tgBot.Bot.Start()
	}()

	defer tgBot.Bot.Stop()
	registerBotHandlers(tgBot)

	bootstrap.BootStrapApp()
}


func registerBotHandlers(tgBot bot.TgBot) {
	commandHandler := chat.NewCommandHandler(&tgBot)
	tgBot.Bot.Handle("/start", commandHandler.StartHandler)
}