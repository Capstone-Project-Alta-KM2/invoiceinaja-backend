package user

type UserFormatter struct {
	Id           int    `json:"id"`
	Fullname     string `json:"nama_lengkap"`
	Email        string `json:"email"`
	NoTlpn       string `json:"no_tlpn"`
	BusinessName string `json:"nama_bisnis"`
	Password     string `json:"password"`
	Avatar       string `json:"avatar"`
	Token        string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		Id:           user.ID,
		Fullname:     user.Fullname,
		Email:        user.Email,
		NoTlpn:       user.NoTlpn,
		BusinessName: user.BusinessName,
		Password:     user.Password,
		Avatar:       user.Avatar,
		Token:        token,
	}

	return formatter
}

type UpdateUserFormatter struct {
	Id           int    `json:"id"`
	Fullname     string `json:"nama_lengkap"`
	Email        string `json:"email"`
	NoTlpn       string `json:"no_tlpn"`
	BusinessName string `json:"nama_bisnis"`
}

func FormatUpdateUser(user User) UpdateUserFormatter {
	formatter := UpdateUserFormatter{
		Id:           user.ID,
		Fullname:     user.Fullname,
		Email:        user.Email,
		NoTlpn:       user.NoTlpn,
		BusinessName: user.BusinessName,
	}

	return formatter
}
