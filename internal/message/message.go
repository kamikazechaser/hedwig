package message

// Message describes incoming API request
type Message struct {
	Service string   `json:"service"`
	To      string   `json:"to"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Params  []string `json:"params"`
}
