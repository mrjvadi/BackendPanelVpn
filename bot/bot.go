package bot

import (
	"fmt"
	event_handler "github.com/mrjvadi/BackendPanelVpn/bot/event-handler"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"gopkg.in/telebot.v4"
)

// Bot represents a messaging bot component.
type Bot struct {
	bus *events.Bus
	bot *telebot.Bot
}

// New creates a new Bot bound to the event bus.
func New(bus *events.Bus, cfg config.BotConfig) *Bot {
	pref := telebot.Settings{
		Token: cfg.Token,
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		fmt.Println("Error creating bot:", err)
	}

	b := &Bot{
		bus: bus,
		bot: bot,
	}

	return b
}

// Start initializes the bot and subscribes to relevant events.
func (b *Bot) Start() {
	// Start the event handler
	busHandler := event_handler.NewBotEventHandler(b.bot, b.bus)
	busHandler.EventStart()

	//b.bus.Subscribe("service:user_created", b.handleUserCreated)
}
