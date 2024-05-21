package models

type NoteRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	UserID  uint64 `json:"userID"`
	//User    User
}
