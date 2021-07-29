package message

type Message struct {
	To string `json:"to"`
	Title string `json:"title"`
	Content string `json:"content"`
}