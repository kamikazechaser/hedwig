package message

// Message represents the common message struct for all services to consume
type Message struct {
	To      string `json:"to"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
