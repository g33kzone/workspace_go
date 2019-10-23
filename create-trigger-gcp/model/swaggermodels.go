package model

//APIMessage represents API message
type APIMessage struct {
	Code int    `json:"code" example:"200"`
	Text string `json:"text" example:"hello world!"`
}
