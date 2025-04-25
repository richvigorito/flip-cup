package types

type Answer struct {
	Message   Message
	Answer string `json:"answer,omitempty"`
}
