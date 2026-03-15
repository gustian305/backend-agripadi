package dto

type UserRegisterRequest struct {
	NamaLengkap     string `json:"nama_lengkap" validate:"required,min=3,max=150"`
	NoTelepon       string `json:"no_telepon" validate:"required,min=10,max=20"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	NoTelepon string `json:"no_telepon" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	ID          string `json:"id"`
	Token       string `json:"token"`
	ExpiresIn   int64  `json:"expires_in"`
	NamaLengkap string `json:"nama_lengkap"`
	NoTelepon   string `json:"no_telepon"`
	Validasi    bool   `json:"validasi"`
}