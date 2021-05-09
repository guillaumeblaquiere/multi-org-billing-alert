package model

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
	} `json:"message"`
}
