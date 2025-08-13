package request

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Email    string `json:"email" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
	Status   string `json:"status" validate:"required,status"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
}
