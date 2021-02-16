package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niggelgame/minecrafter/pkg/commands"
	"github.com/niggelgame/minecrafter/pkg/minecraft"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Token      string
	Prefix     string
	Connection *minecraft.Connection
}

type Bot struct {
	session         *discordgo.Session
	commandHandler  *commands.CommandHandler
	commandRegistry *commands.CommandRegistry
	connection      *minecraft.Connection
	config          *Config
}

func New(conf Config) *Bot {
	dg, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		zap.L().Fatal("Failed to create Discord session.", zap.Error(err))
	}

	commandRegistry := commands.NewCommandRegistry()
	commandHandler := commands.NewHandler(commandRegistry, conf.Prefix)

	b := &Bot{
		session:         dg,
		commandHandler:  commandHandler,
		commandRegistry: commandRegistry,
		connection:      conf.Connection,
		config:          &conf,
	}

	b.registerHandler()
	b.registerCommands()

	return b
}

func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		zap.L().Fatal("Error opening Connection.", zap.Error(err))
		return err
	}

	err = b.session.UpdateListeningStatus(fmt.Sprintf("%shelp", b.config.Prefix))
	if err != nil {

	}

	// Wait here until CTRL-C or other term signal is received.
	zap.L().Info("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = b.session.Close()
	if err != nil {
		zap.L().Fatal("Failed to close down connection cleanly.", zap.Error(err))

		return err
	}
	return nil
}

func (b *Bot) registerHandler() {
	b.session.AddHandler(b.commandHandler.HandleMessageEvent)
}

func (b *Bot) registerCommands() {
	b.commandRegistry.RegisterCommand(commands.NewMinecraftCommand(b.connection))
	b.commandRegistry.RegisterCommand(commands.NewWhitelistCommand(b.connection))
	b.commandRegistry.RegisterCommand(commands.NewHelpCommand())
}
