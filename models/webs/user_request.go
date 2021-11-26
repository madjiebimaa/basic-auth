package webs

type UserRegisterRequest struct {
	Name     string `validate:"required,min=3" json:"name"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=5" json:"password"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=5" json:"password"`
}
