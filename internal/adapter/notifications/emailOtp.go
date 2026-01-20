package notifications

import (
	"net/smtp"
	"fmt"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type emailNotifier struct {
	host     string
	port     string
	username string
	password string
	from 	 string
}

func NewemailNodifier(h, p, u, pass string) domain.NotificationClint{
	return &emailNotifier{
		host: h,
		port: p,
		username: u,
		password: pass,
	}
}
func (e *emailNotifier) SendOtp(toEmail,code string )error{
	auth:=smtp.PlainAuth("",e.username,e.password,e.host)
		msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Verification Code for the sound core\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"Your code is: %s\r\n", toEmail, code))

	addr:=e.host+":"+e.port
	return smtp.SendMail(addr, auth, e.username, []string{toEmail}, msg)

}