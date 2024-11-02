package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/repository"
)

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessToken)
}
func (b *Bot) initAutorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(b.messages.Start, authLink))
	msg.ReplyToMessageID = message.MessageID
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestToken); err != nil {
		return "", err
	}
	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}
func (b *Bot) generateRedirectURL(chatID int64) string {
	url := fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
	return url
}
