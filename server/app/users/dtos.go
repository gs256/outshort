package users

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
