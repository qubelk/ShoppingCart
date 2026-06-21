package user

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type LoginResponse struct {
	User User `json:"user"`
}
