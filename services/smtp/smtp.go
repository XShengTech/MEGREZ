package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"megrez/services/config"
	"net/mail"
	"net/smtp"
	"strconv"
)

const emailTitle = "MEGREZ 天权算能聚联计算平台"

func Send(toAddr, title string, html string) error {
	smtpConf := config.GetSmtp()

	if smtpConf.Host == "" || smtpConf.Port == 0 || smtpConf.User == "" || smtpConf.Password == "" {
		return errors.New("SMTP smtpConfiguration is not set")
	}

	if !smtpConf.SSL {
		auth := smtp.PlainAuth("", smtpConf.User, smtpConf.Password, smtpConf.Host)
		msg := append([]byte("Subject: "+title+" - "+emailTitle+" \r\n"+
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"),
			html...)
		err := smtp.SendMail(smtpConf.Host+":"+strconv.Itoa(smtpConf.Port), auth, smtpConf.User, []string{toAddr}, msg)
		if err != nil {
			return errors.New("fail to send email, Error:" + err.Error())
		}
		return nil
	}

	from := mail.Address{Name: "MEGREZ", Address: smtpConf.User}
	to := mail.Address{Name: "", Address: toAddr}
	auth := smtp.PlainAuth("", smtpConf.User, smtpConf.Password, smtpConf.Host)
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = title + " - " + emailTitle
	headers["Content-Type"] = "text/html;charset=UTF-8"
	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + html
	tlssmtpConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpConf.Host,
	}
	conn, err := tls.Dial("tcp", smtpConf.Host+":"+strconv.Itoa(smtpConf.Port), tlssmtpConfig)
	if err != nil {
		return errors.New("fail to connect to the server, Error:" + err.Error())
	}
	c, err := smtp.NewClient(conn, smtpConf.Host)
	if err != nil {
		return errors.New("fail to create smtp client, Error:" + err.Error())
	}
	if err = c.Auth(auth); err != nil {
		return errors.New("fail to auth, Error:" + err.Error())
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return errors.New("fail to set from address, Error:" + err.Error())
	}

	if err = c.Rcpt(to.Address); err != nil {
		return errors.New("fail to set to address, Error:" + err.Error())
	}
	w, err := c.Data()
	if err != nil {
		return errors.New("fail to get smtp data writer, Error:" + err.Error())
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return errors.New("fail to write data, Error:" + err.Error())
	}

	err = w.Close()
	if err != nil {
		return errors.New("fail to close smtp data writer, Error:" + err.Error())
	}

	c.Quit()

	return nil
}
