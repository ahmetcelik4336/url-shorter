package dto

type LoginRequest struct {
	Email    string `form:"email" validate:"required,email" json:"email"`
	Password string `form:"password" validate:"required,min=5" json:"password"`
}

type RegisterRequest struct {
	Email           string `form:"email" validate:"required,email" json:"email"`
	Password        string `form:"password" validate:"required,min=5" json:"password"`
	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password" json:"password_confirm"`
	Name            string `form:"name" validate:"required" json:"name"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
