package request

type Login struct {
	Email    string `json:"email"` // Email or Username
	Password string `json:"password"`
}

type GetPosts struct {
	StartId int    `json:"startId"`
	Session string `json:"session"`
}
