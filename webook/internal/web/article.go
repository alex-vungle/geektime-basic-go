package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService,
	l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/articles")
	// 在有 list 等路由的时候，无法这样注册
	// g.GET("/:id", a.Detail)
	g.GET("/detail/:id", a.Detail)
	// 理论上来说应该用 GET的，但是我实在不耐烦处理类型转化
	// 直接 POST，JSON 转一了百了。
	g.POST("/list", a.List)

	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)

	// 准备给读者用的服务，但是暂时没太大的必要实现
	//pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	//pub.GET("/pub/:id", a.PubDetail)
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	if err := a.svc.Withdraw(ctx, usr.Id, req.Id); err != nil {
		a.l.Error("设置为尽自己可见失败", logger.Error(err),
			logger.Field{Key: "id", Value: req.Id})
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (a *ArticleHandler) List(ctx *gin.Context) {
	type Req struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}

	// 对于批量接口来说，要小心批次大小
	if req.Limit > 100 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 我会倾向于不告诉前端批次太大
			// 因为一般你和前端一起完成任务的时候
			// 你们是协商好了的，所以会进来这个分支
			// 就表明是有人跟你过不去
			Msg: "请求有误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}

	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	arts, err := a.svc.List(ctx, usr.Id, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: slice.Map[domain.Article, ArticleVo](arts,
			func(idx int, src domain.Article) ArticleVo {
				return ArticleVo{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					// 这个列表请求，不需要返回内容
					//Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author
					Ctime: src.Ctime.Format(time.DateTime),
					Utime: src.Utime.Format(time.DateTime),
				}
			}),
	})
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	art, err := a.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		a.l.Error("非法访问文章，创作者 ID 不匹配",
			logger.Int64("uid", usr.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Publish(ctx, req.toDomain(usr.Id))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Save(ctx, req.toDomain(usr.Id))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存数据失败", logger.Field{Key: "error", Value: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}
