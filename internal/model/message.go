package model

type Message struct {
	IsInitialized bool   `json:"isInitialized"`
	SourceMacAddr string `json:"sourceMacAddr"`
}
