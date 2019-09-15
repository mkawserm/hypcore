package gqltypes

type ErrorType struct {
	Messages    []interface{} `json:"messages"`
	MessageType string        `json:"type"`
	Group       int           `json:"group"`
	Code        int           `json:"code"`
	TraceID     string        `json:"trace_id"`
}

func NewErrorType() *ErrorType {
	return &ErrorType{Messages: make([]interface{}, 0)}
}

func (e *ErrorType) AddMessage(message interface{}) {
	e.Messages = append(e.Messages, message)
}

func (e *ErrorType) AddStringMessage(message string) {
	e.Messages = append(e.Messages, map[string]string{"message": message})
}
