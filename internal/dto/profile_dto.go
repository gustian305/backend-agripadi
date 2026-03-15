package dto

type UpdateProfileRequest struct {
	NamaLengkap     string `json:"nama_lengkap"`
	NoTelepon       string `json:"no_telepon"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdateProfileResponse struct {
	ID          string `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	NoTelepon   string `json:"no_telepon"`
}