package user

type InputRegister struct {
	Fullname     string `json:"fullname" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	NoTlpn       string `json:"phone_number" binding:"required"`
	BusinessName string `json:"company" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type InputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InputCheckEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type InputUpdate struct {
	Fullname     string `json:"fullname" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	NoTlpn       string `json:"phone_number" binding:"required"`
	BusinessName string `json:"company" binding:"required"`
}

type InputChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
