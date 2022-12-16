package models

type Response struct {
	Status       string      `json:"status"`
	InnerMessage string      `json:"inner_message"`
	Message      string      `json:"message"`
	Body         interface{} `json:"body,omitempty"`
}