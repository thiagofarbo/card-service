package models

type CardRequest struct {
	Number string `json:"number"`
	UserID uint64 `json:"userID"`
}
