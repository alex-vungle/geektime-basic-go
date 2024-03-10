package domain

type User struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
}
