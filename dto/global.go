package dto

type Message struct {
	Message string `json:"message"`
}

type ValidationError struct {
	FailedField string
	Tag         string
	Value       string
}
