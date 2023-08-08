package web

import (
	"encoding/json"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/webook-server/internal/domain"
	"webook/webook-server/internal/service"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp2.Regexp
	passwordExp *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	passwordExp := regexp2.MustCompile(passwordRegexPattern, regexp2.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (uh *UserHandler) RegisterRouters(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", uh.SignUp)
	ug.POST("/login", uh.Login)
	ug.POST("/edit", uh.Edit)
	ug.GET("/profile", uh.Profile)
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	//Bind方法会根据 Content-Type 来解析你的数据到req里面
	//解析错了，直接返回400错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := uh.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
	}

	ok, err = uh.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，需要包含数字、字母和特殊字符")
		return
	}

	//调用 svc 的方法
	err = uh.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	//下面放数据库操作
}

func (uh *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := uh.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}

	//步骤2
	//登录成功设置session
	sess := sessions.Default(ctx)
	//可以随便设置你放在session中的值
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}

func (uh *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
		Avatar       string `json:"avatar"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if len(req.Nickname) >= 12 {
		ctx.String(http.StatusOK, "昵称超过长度限制")
		return
	}

	if len(req.Introduction) >= 200 {
		ctx.String(http.StatusOK, "个人简介超出字数限制")
		return
	}

	//_, err := time.Parse("1962-01-01", req.Birthday)
	//if err != nil {
	//	ctx.String(http.StatusOK, "生日日期格式不对")
	//	return
	//}

	sess := sessions.Default(ctx)
	userId := sess.Get("userId").(int64)
	err := uh.svc.Edit(ctx, domain.User{
		Id:           userId,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
		Avatar:       req.Avatar,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "编辑成功")

}

func (uh *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId").(int64)
	u, err := uh.svc.Profile(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	res, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	ctx.String(http.StatusOK, string(res))
}
