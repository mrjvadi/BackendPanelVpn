package event_handler

import (
	"github.com/mrjvadi/BackendPanelVpn/events"
	"gopkg.in/telebot.v4"
)

// BotEventHandler handles events related to the bot.
type BotEventHandler struct {
	bot *telebot.Bot
	bus *events.Bus
}

// NewBotEventHandler creates a new BotEventHandler instance.
func NewBotEventHandler(bot *telebot.Bot, bus *events.Bus) *BotEventHandler {
	return &BotEventHandler{
		bot: bot,
		bus: bus,
	}
}

// EventStart initializes the event handler for the bot.
func (b *BotEventHandler) EventStart() {
	b.bus.Subscribe("any:login", b.eventLogin)
}
