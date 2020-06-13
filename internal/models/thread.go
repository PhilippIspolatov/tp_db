package models

import "time"

type Thread struct {
	Author string `json:"author"`
	Created time.Time `json:"created"`
	Forum string `json:"forum"`
	Id uint64 `json:"id"`
	Message string `json:"message"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	Votes int64 `json:"votes"`
}
