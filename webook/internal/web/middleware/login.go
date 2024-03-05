package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	ignorePaths map[string]string
	abortStatus int
}

//func (m *LoginMiddlewareBuilder) BuildLogin() gin.HandlerFunc {
//
//}

func (m *LoginMiddlewareBuilder) AddIgnorePath(path ...string) *LoginMiddlewareBuilder {
	for _, p := range path {
		m.ignorePaths[p] = p
	}
	return m
}

func (m *LoginMiddlewareBuilder) AbortStatus(status int) *LoginMiddlewareBuilder {
	m.abortStatus = status
	return m
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册一下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		//_, ok := m.ignorePaths[path]
		//if ok {
		//	return
		//}

		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			//ctx.AbortWithStatus(m.abortStatus)
			return
		}

		now := time.Now()
		//ctx.Next()// 执行业务
		// 在执行业务之后搞点什么
		//duration := time.Now().Sub(now)

		// 我怎么知道，要刷新了呢？
		// 假如说，我们的策略是每分钟刷一次，我怎么知道，已经过了一分钟？
		const updateTimeKey = "update_time"
		// 试着拿出上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 你这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}

// 现在我说，我要增加更多的忽略字段
// 加了一个参数，是不兼容的修改
// 你只能加字段

var (
	IgnorePath  map[string]string
	AbortStatus int
)

func MultipleServer() {
	// 业务 server
	// 我要求 server1 返回 401
	server1 := gin.Default()
	AbortStatus = 401
	server1.Use(CheckLoginV1())
	// 管理后台 server
	// server2 返回 403
	AbortStatus = 403
	server2 := gin.Default()
	server2.Use(CheckLoginV1())

	// 最终两个 server 都是返回 403
}

func CheckLoginV1() gin.HandlerFunc {
	// 注册一下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path

		_, ok := IgnorePath[path]
		if ok {
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(AbortStatus)
			return
		}

		now := time.Now()
		//ctx.Next()// 执行业务
		// 在执行业务之后搞点什么
		//duration := time.Now().Sub(now)

		// 我怎么知道，要刷新了呢？
		// 假如说，我们的策略是每分钟刷一次，我怎么知道，已经过了一分钟？
		const updateTimeKey = "update_time"
		// 试着拿出上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 你这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}

func CheckLogin(ignorePaths map[string]string, status int) gin.HandlerFunc {
	// 注册一下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path

		_, ok := ignorePaths[path]
		if ok {
			return
		}

		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		//ctx.Next()// 执行业务
		// 在执行业务之后搞点什么
		//duration := time.Now().Sub(now)

		// 我怎么知道，要刷新了呢？
		// 假如说，我们的策略是每分钟刷一次，我怎么知道，已经过了一分钟？
		const updateTimeKey = "update_time"
		// 试着拿出上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 你这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}
