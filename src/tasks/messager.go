package tasks

import (
	"fmt"
	"log"

	"github.com/imroc/req/v3"

	"desprit/hammertime/src/config"
)

type Messager interface {
	SendMessage(message string) error
}

type TelegramMessager struct{}

func NewTelegramMessager() *TelegramMessager {
	return &TelegramMessager{}
}

func (t *TelegramMessager) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.GetConfig().TG_BOT_TOKEN)
	data := map[string]string{"chat_id": config.GetConfig().TG_CHAT_ID, "text": message}
	resp, err := req.R().SetFormData(data).Post(url)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return err
	}
	log.Printf("Notification sent: %v", resp)
	return nil
}
