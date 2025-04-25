package types

type Message struct {
  Type   string      `json:"type"`
  ID     string      `json:"id,omitempty"`
  Name   string      `json:"name,omitempty"`
  Answer interface{} `json:"answer,omitempty"`
}

