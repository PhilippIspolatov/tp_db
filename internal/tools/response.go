package tools

type Body map[string]interface{}

type Error struct {
	ErrMsg string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}