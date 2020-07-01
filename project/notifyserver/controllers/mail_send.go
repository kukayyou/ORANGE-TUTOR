package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/kukayyou/commonlib/mylog"
	"math/rand"
	"notifyserver/config"
	"time"
)

const (
	SOURCE_EMAIL_ADDR   = "zhangkaikay@126.com"
	SOURCE_EMAIL_HOST   = "smtp.126.com"
	SOURCE_EMAIL_PASSWD = "!zhangkai1025"
	SOURCE_EMAIL_PORT   = 465
)

type SendMailController struct {
	BaseController
}

type SendMailRequest struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"` //token
}

func (this SendMailController) SendMailApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params SendMailRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	if err := this.CheckToken(params.UserID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		return
	}
	captcha := creatCaptcha(params.UserID)
	mailBody := fmt.Sprintf("【橘子科技】验证码%s，5分钟内有效。验证码提供给他人可能导致账号被盗，请勿泄漏，谨防被骗。", captcha)
	if err := sendMail(params.Email, mailBody); err != nil {
		this.Resp.Code = SEND_MAIL_ERROR
		this.Resp.Msg = "send mail failed!"
		return
	}
	insertCaptchaIntoRedis(this.GetRequestId(), captcha, params.UserID)
	//insertDataToMongo()
}

func sendMail(targetAddr, mailBody string) error {
	mail := gomail.NewMessage()
	mail.SetAddressHeader("From", SOURCE_EMAIL_ADDR, "橘子科技") // 发件人
	mail.SetHeader("To",                                     // 收件人
		mail.FormatAddress(targetAddr, targetAddr),
	)
	mail.SetHeader("Subject", "橘子科技")   // 主题
	mail.SetBody("text/html", mailBody) // 正文

	d := gomail.NewDialer(SOURCE_EMAIL_HOST, SOURCE_EMAIL_PORT, SOURCE_EMAIL_ADDR, SOURCE_EMAIL_PASSWD) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(mail); err != nil {
		mylog.Error("send mail failed, targetAddr:%s,error:%s", targetAddr, err.Error())
		return err
	}
	return nil
}

//生成随机验证码
func creatCaptcha(UserID int64) (captcha string) {
	captcha = fmt.Sprintf("%d%d%d%d", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9))
	return
}

func insertCaptchaIntoRedis(requestID, captcha string, userID int64) {
	cachekey := fmt.Sprintf("%d", userID)

	err := config.RedisClient.Set(cachekey, captcha, time.Second*60*5).Err()
	if err != nil {
		mylog.Error("requestId:%s, cachekey:%s, insert redis error:%s", requestID, cachekey, err.Error())
	} else {
		mylog.Info("requestId:%s, cachekey:%s, data:%s", requestID, cachekey, captcha)
	}
}

func checkCaptcha(requestID, captcha string, userID int64) error {
	cachekey := fmt.Sprintf("%d", userID)

	data := config.RedisClient.Get(cachekey).Val()
	if len(data) <= 0 {
		mylog.Error("requestId:%s, checkCaptcha error, captcha is time out", requestID, cachekey)
		return fmt.Errorf("check captcha error, captcha is time out")
	} else if data != captcha {
		mylog.Error("requestId:%s, checkCaptcha error, captcha is not equal", requestID, cachekey)
		return fmt.Errorf("check captcha error, captcha is not equal")
	}
	mylog.Info("requestId:%s, check captcha success cachekey:%s, data:%s", requestID, cachekey)
	return nil
}

func insertDataToMongo() {
	collection := config.Client.Database("orange_tutor").Collection("test")
	type Attendee struct {
		EventID    int64   `json:"eventId"`
		HostID     int64   `json:"hostId"`
		AttendList []int64 `json:"userId"`
	}

	attendee := make([]interface{}, 0)
	for i := 0; i < 1000000; i++ {
		attendList := make([]int64, 0)
		count := rand.Intn(100)
		for m := 0; m <= count; m++ {
			attendList = append(attendList, int64(m))
		}
		data := Attendee{
			EventID:    int64(100 + i),
			HostID:     int64(i % 100),
			AttendList: attendList,
		}
		attendee = append(attendee, data)
	}
	collection.InsertMany(context.TODO(), attendee)
}
