package GetMessage

import (
	"Message/Models"
	"rpc/Message"
)

var OnceMaxMessageNumber = 10

func RefreshText(sender, receiver string) []*Message.Text {
	Messages := Models.GetMessages[*Models.Text](sender, receiver, OnceMaxMessageNumber)
	texts := make([]*Message.Text, 0, len(Messages))
	for _, message := range Messages {
		var text = Message.Text{}
		text.Content = string(message.GetContent())
		text.SenderAccountNum = message.SenderAccountNum
		text.ReceiverAccountNum = message.ReceiverAccountNum
		text.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		texts = append(texts, &text)
	}

	return texts
}
