package web

import (
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

var _ handler = (*OAuth2WechatHandler)(nil)

type OAuth2WechatHandler struct {
	// 这边也可以直接定义成 wechat.Service
	// 但是为了保持使用 mock 来测试，这里还是用了接口
	wechatSvc       wechat.Service
	userSvc         service.UserService
	stateCookieName string
	jwtHandler
}

func NewOAuth2WechatHandler(service wechat.Service,
	userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatSvc: service,
		userSvc:   userSvc,
		// 万一后续我们要改，也可以做成可配置的。
		stateCookieName: "jwt-state",
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.OAuth2URL)
	// 这边用 Any 万无一失
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证 state
	err := h.verifyState(ctx)
	if err != nil {
		// 实际上，但凡进来这里，就说明有人搞你，
		// 因此这边要做好监控和告警
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常，请重试",
		})
		return
	}

	code := ctx.Query("code")
	info, err := h.wechatSvc.VerifyCode(ctx, code)
	if err != nil {
		// 实际上这个错误，也有可能是 code 不对
		// 但是给前端的信息没有太大的必要区分究竟是代码不对还是系统本身有问题
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里就是登录成功
	// 所以你需要设置 JWT
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		// 基本上，如果进来这里，就可以认为是有人在搞鬼。
		return fmt.Errorf("%w, 无法获得 cookie", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil {
		return fmt.Errorf("%w, cookie 不是合法 JWT token", err)
	}
	if sc.State != state {
		return errors.New("state 被篡改了")
	}
	return nil
}

func (h *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.wechatSvc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误，请稍后再试",
		})
		return
	}
	err = h.setStateCookie(ctx, state)
	if err != nil {
		// 理论上你也可以考虑忽略这个错误，不影响扫码登录
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误，请稍后再试",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
	return
}

// setStateCookie 只有微信这里用，所以定义在这里
func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
	})
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr,
		600,
		// 限制在只能在这里生效。
		"/oauth2/wechat/callback",
		// 这边把 HTTPS 协议禁止了。不过在生产环境中要开启。
		"", false, true)
	return nil
}
