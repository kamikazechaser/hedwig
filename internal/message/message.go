package message

// Message represents the common message struct for all services to consume
type Message struct {
	Service string `json:"service"`
	To      string `json:"to"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
