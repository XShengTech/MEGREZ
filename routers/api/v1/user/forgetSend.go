package user

import (
	"fmt"
	"megrez/libs/crypto"
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/config"
	"megrez/services/database"
	"megrez/services/redis"
	"megrez/services/smtp"
	"time"

	"github.com/kataras/iris/v12"
)

type forgetSendStruct struct {
	Email string `json:"email"`
}

const forgetPasswordRedisKeyPrefix = "forget:user:"
const forgetPasswordUrlPrefix = "/reset/"
const forgetPasswordTitle = "重置密码"
const forgetPasswordHTMLFormat = `
<div>
    <table cellpadding="0" align="center" style="overflow:hidden;background:#fff;margin:0 auto;text-align:left;position:relative;font-size:14px; font-family:'lucida Grande',Verdana;line-height:1.5;box-shadow:0 0 3px #ccc;border:1px solid #ccc;border-radius:5px;border-collapse:collapse;">
        <tbody>
            <tr>
                <th valign="middle" style="height:38px;color:#fff; font-size:14px;line-height:38px; font-weight:bold;text-align:left;padding:10px 24px 6px; border-bottom:1px solid #467ec3;background:#518bcb;border-radius:5px 5px 0 0;">
                    MEGREZ 天权算能聚联计算平台</th>
            </tr>
            <tr>
                <td>
                    <div style="padding:20px 35px 40px;">
                        <h2 style="font-weight:bold;margin-bottom:5px;font-size:14px;">Hello, %s:</h2>
                        <p style="margin-top:20px">
                            请在15分钟内点击链接： <a href="%s">%s</a> &nbsp;进行密码重置操作，十五分钟后该链接将会失效.
                        </p>
                        <p style="margin-top:20px">
                            为了保护你的账户,请不要使用单一的密码来进行重置。
                        </p>
                        <p style="margin-top:20px">
                            如果您有任何问题，请联系系统管理员以获得更多信息与支持。
                        </p>
                        <p style="margin-left:2em;"></p>
                        <p style="text-indent:0;text-align:right;">MEGREZ 天权算能聚联计算平台</p>
                    </div>
                </td>
            </tr>

        </tbody>
    </table>
</div>
`

func forgetSendHandler(ctx iris.Context) {
	var req forgetSendStruct
	if err := ctx.ReadJSON(&req); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.Email == "" {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	user := models.Users{
		Email: req.Email,
	}
	result := database.DB.Where(&user).First(&user)
	if result.Error != nil {
		l.Error("get user error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeUserNotExist, iris.StatusInternalServerError)
		return
	}

	rdb := redis.RawDB

	forgetUrl := crypto.Hex(32)
	err := rdb.Set(ctx, forgetPasswordRedisKeyPrefix+forgetUrl, user.ID, 15*time.Minute).Err()
	if err != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("Set Redis error: %v", err)
		return
	}

	forgetUrl = config.GetSystemBaseUrl() + forgetPasswordUrlPrefix + forgetUrl
	err = smtp.Send(user.Email, forgetPasswordTitle, fmt.Sprintf(forgetPasswordHTMLFormat, user.Username, forgetUrl, forgetUrl))
	if err != nil {
		middleware.Error(ctx, middleware.CodeServeBusy, iris.StatusInternalServerError)
		l.Error("Send SMTP Error: %v", err)
		return
	}

	middleware.Success(ctx)
}
