package models

type Vote struct {
	Nickname string `json:"nickname"`
	Thread uint64 `json:"thread"`
	Voice int64 `json:"voice"`
}
