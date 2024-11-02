package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errorInvalidLink        = errors.New("invalid link")
	errorUnautorized        = errors.New("user is not autorized")
	errorUnabletoSave       = errors.New("unable to save")
	errorUnavailableCommand = errors.New("unavailable command")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)
	switch err {
	case errorInvalidLink:
		msg.Text = b.messages.InvalidURL
		b.bot.Send(msg)
	case errorUnautorized:
		msg.Text = b.messages.Unautorized
		b.bot.Send(msg)
	case errorUnabletoSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	case errorUnavailableCommand:
		msg.Text = b.messages.UnknownCommand
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
