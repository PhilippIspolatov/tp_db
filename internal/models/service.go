package models

type Service struct {
	Forum uint64 `json:"forum"`
	Thread uint64 `json:"thread"`
	User uint64 `json:"user"`
	Post uint64 `json:"post"`
}
