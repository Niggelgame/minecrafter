package main

import (
	"github.com/niggelgame/minecrafter/pkg/bot"
	"github.com/niggelgame/minecrafter/pkg/config"
	"github.com/niggelgame/minecrafter/pkg/minecraft"
	"go.uber.org/zap"
)

// Project structure heavily inspired by https://github.com/ForYaSee/NindoBot

func main() {
	conf := config.Load()

	if conf.IsDev {
		logger, err := zap.NewDevelopment()
		if err == nil {
			zap.ReplaceGlobals(logger)
		}
	}

	mC, err := minecraft.NewConnection(minecraft.Config{
		Address:  conf.RconAddress,
		Password: conf.RconPassword,
	})

	if err != nil {
		zap.L().Warn("Could not connect to Minecraft RCON.", zap.Error(err))
	}

	b := bot.New(bot.Config{
		Token:      conf.BotToken,
		Prefix:     conf.BotPrefix,
		Connection: mC,
	})

	// go func() {
	zap.L().Info("Bot started.")
	zap.L().Fatal("Error while serving bot.", zap.Error(b.Start()))
	// }()

}
