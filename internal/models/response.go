package models

type APIResponse struct {
	Error    *APIError              `json:"error,omitempty"`
	Response map[string]interface{} `json:"response,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
}

type APIError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
