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

func NewemailNodifier(h, p, u, pass,f string) domain.NotificationClint{
	return &emailNotifier{
		host: h,
		port: p,
		username: u,
		password: pass,
		from:f,
	}
}
func (e *emailNotifier) SendOtp(toEmail,code string )error{
	auth:=smtp.PlainAuth("",e.username,e.password,e.host)

	subject := "Subject: Verify your account\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h3>Welcome to My Ecommerce!</h3>
				<p>Your verification code is: <strong>%s</strong></p>
				<p>This code expires in 10 minutes.</p>
			</body>
		</html>
	`, code)
	
	msg := []byte(subject + mime + body)

	addr := e.host + ":" + e.port
	err := smtp.SendMail(addr, auth, e.from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}