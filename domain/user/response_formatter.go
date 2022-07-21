package user

type UserFormatter struct {
	Id          int    `json:"id"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Company     string `json:"company"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	Token       string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		Id:          user.ID,
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.NoTlpn,
		Company:     user.BusinessName,
		Password:    user.Password,
		Avatar:      user.Avatar,
		Token:       token,
	}

	return formatter
}

type UpdateUserFormatter struct {
	Id          int    `json:"id"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Company     string `json:"company"`
}

func FormatUpdateUser(user User) UpdateUserFormatter {
	formatter := UpdateUserFormatter{
		Id:          user.ID,
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.NoTlpn,
		Company:     user.BusinessName,
	}

	return formatter
}
