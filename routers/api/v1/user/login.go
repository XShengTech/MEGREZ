package user

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/routers/api/v1/sessions"
	"megrez/services/database"
	"megrez/services/redis"
	"time"

	"github.com/kataras/iris/v12"
)

type userLoginStruct struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	// HCaptchaToken string
	// CaptchaType   string
	// Key           string
}

func loginHandler(ctx iris.Context) {
	sess := sessions.Session()
	session := sess.Start(ctx)
	sess.UseDatabase(redis.DB)

	var userReq userLoginStruct
	if err := ctx.ReadJSON(&userReq); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if userReq.Account == "" || userReq.Password == "" {
		middleware.Error(ctx, middleware.CodeLoginError, iris.StatusBadRequest)
		return
	}

	user := models.Users{}
	result := database.DB.Where("username = ?", userReq.Account).Or("email = ?", userReq.Account).First(&user)
	if result.Error != nil {
		middleware.Error(ctx, middleware.CodeLoginError, iris.StatusBadRequest)
		return
	}

	if !user.CheckPassword(userReq.Password) {
		middleware.Error(ctx, middleware.CodeLoginError, iris.StatusBadRequest)
		return
	}

	if user.Role < 1 {
		middleware.Error(ctx, middleware.CodeForbidden, iris.StatusForbidden)
		return
	}

	session.Set("authenticated", true)
	session.Set("userId", user.ID)
	session.Set("username", user.Username)
	session.Set("email", user.Email)
	session.Set("role", user.Role)
	session.Set("loginTime", time.Now().Unix())

	profile := profile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Balance:  user.Balance,
	}

	middleware.Result(ctx, profile)
}
