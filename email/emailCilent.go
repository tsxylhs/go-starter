package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"strconv"
	"time"

	code "github.com/tsxylhs/go-starter/domain"
)

func SendMail(emailConfig code.EmailConf, toUsers []string, data *interface{}) error {
	e := NewEmail()

	e.From = emailConfig.UserName
	e.To = toUsers

	e.Subject = "计划统计"
	//emailModel.Title = e.Subject
	t, err := template.ParseFiles("template.html")
	if err != nil {

		return err
	}
	startTime, _ := time.ParseInLocation("15:04", "9:00", time.Local)
	start := startTime.AddDate(time.Now().Year(), int(time.Now().Month()-1), time.Now().Day()-2).Format("2006-01-02 15:04")
	endtime := startTime.AddDate(time.Now().Year(), int(time.Now().Month()-1), time.Now().Day()-1).Format("2006-01-02 15:04")
	body := new(bytes.Buffer)
	t.Execute(body, struct {
		Message1 string
		Message2 string
		// JigStatistics    []model.JigStatisticsDto
		// MaintainPlanData []model.MaintenancePlanDto
	}{
		Message1: "使用情况总计 统计时间：" + start + "--" + endtime,
		Message2: "test    统计时间：" + endtime,
		// JigStatistics:    data.JigStatisticsDto,
		// MaintainPlanData: data.MaintainPlanData,
	})
	e.HTML = body.Bytes()
	e.AttachFile("./test.xls")
	auth := NewLoginAuth(emailConfig.UserName, emailConfig.Password)
	err = e.Send(emailConfig.Host+":"+strconv.Itoa(emailConfig.Port), auth)
	if err != nil {
		//log.Logger.Logger.Error("邮件发送失败", zap.Error(err))
		return nil
	} else {
		//log.Logger.Logger.Info("邮件发送成功！")
	}
	return nil

}

type LoginAuth struct {
	username, password string
}

func NewLoginAuth(username, password string) smtp.Auth {
	return &LoginAuth{username, password}
}

func (a *LoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *LoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}

func SendMailVo(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	host, _, _ := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("call dial")
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host, InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			fmt.Println("call start tls")
			return err
		}
	}

	// if a != nil {
	// 	if ok, _ := c.Extension("AUTH"); ok {
	// 		if err = c.Auth(a); err != nil {
	// 			fmt.Println("check auth with err:", err)
	// 			return err
	// 		}
	// 	}
	// }
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// func main() {
// 	auth := NewLoginAuth("username", "password")

// 	to := []string{"收件人邮箱"}
// 	msg := []byte("这是一封来自go的测试邮件")
// 	err := SendMail("smtphostwithport", auth, "发件人邮箱", to, msg)
// 	if err != nil {
// 		fmt.Println("with err:", err)
// 	}
// 	fmt.Println("please check mailbox")
// }
