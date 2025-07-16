package Models

type WebSocketData struct {
	Data       string `json:"data"`
	AccountNum string `json:"accountNum"`
	Receiver   string `json:"receiver"`
	Type       string `json:"type"`
	CreatedAt  string `json:"createdAt"`
}
