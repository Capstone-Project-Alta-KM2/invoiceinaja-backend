package utils

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

func SendMail(destination string, newPass string) {
	abc := gomail.NewMessage()

	abc.SetHeader("From", "bkppbogor12@gmail.com")
	abc.SetHeader("To", destination)
	abc.SetHeader("Subject", "Reset Password DKPP Kab.Bogor")
	abc.SetBody("text/html", "Berikut ini kami berikan kata sandi untuk akun anda. <br/> <b>"+newPass+"</b> <br/> Mohon untuk segera login dan ganti password baru.")

	a := gomail.NewDialer("smtp.gmail.com", 587, "bkppbogor12@gmail.com", "kokolopo")

	if err := a.DialAndSend(abc); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
