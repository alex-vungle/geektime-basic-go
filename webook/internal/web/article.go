package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     service.ArticleService
	intrSvc service.InteractiveService
	l       logger.LoggerV1
	biz     string
}

func NewArticleHandler(l logger.LoggerV1,
	svc service.ArticleService,
	intrSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		l:       l,
		svc:     svc,
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	//g.PUT("/", h.Edit)
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)

	g.GET("/trends", h.TopLikes)

	// 创作者接口
	g.GET("/detail/:id", h.Detail)
	// 按照道理来说，这边就是 GET 方法
	// /list?offset=?&limit=?
	g.POST("/list", h.List)

	pub := g.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	// 传入一个参数，true 就是点赞, false 就是不点赞
	pub.POST("/like", h.Like)
	pub.POST("/collect", h.Collect)
}

// Edit 接收 Article 输入，返回一个 ID，文章的 ID
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      int64
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		h.l.Error("保存文章数据失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//val, ok := ctx.Get("user")
	//if !ok {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 4,
	//		Msg:  "未登录",
	//	})
	//	return
	//}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("发表文章失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	err := h.svc.Withdraw(ctx, uc.Uid, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("撤回文章失败",
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}
	// 我要不要检测一下？
	uc := ctx.MustGet("user").(jwt.UserClaims)
	arts, err := h.svc.GetByAuthor(ctx, uc.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("查找文章列表失败",
			logger.Error(err),
			logger.Int("offset", page.Offset),
			logger.Int("limit", page.Limit),
			logger.Int64("uid", uc.Uid))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: slice.Map[domain.Article, ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
			return ArticleVo{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),

				//Content:  src.Content,
				AuthorId: src.Author.Id,
				// 列表，你不需要
				Status: src.Status.ToUint8(),
				Ctime:  src.Ctime.Format(time.DateTime),
				Utime:  src.Utime.Format(time.DateTime),
			}
		}),
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.l.Warn("查询文章失败，id 格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}
	art, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("查询文章失败",
			logger.Int64("id", id),
			logger.Error(err))
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if art.Author.Id != uc.Uid {
		// 有人在搞鬼
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("非法查询文章",
			logger.Int64("id", id),
			logger.Int64("uid", uc.Uid))
		return
	}

	vo := ArticleVo{
		Id:    art.Id,
		Title: art.Title,
		//Abstract: art.Abstract(),

		Content:  art.Content,
		AuthorId: art.Author.Id,
		// 列表，你不需要
		Status: art.Status.ToUint8(),
		Ctime:  art.Ctime.Format(time.DateTime),
		Utime:  art.Utime.Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, Result{Data: vo})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.l.Warn("查询文章失败，id 格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}

	var (
		eg   errgroup.Group
		art  domain.Article
		intr domain.Interactive
	)

	uc := ctx.MustGet("user").(jwt.UserClaims)
	eg.Go(func() error {
		var er error
		art, er = h.svc.GetPubById(ctx, id, uc.Uid)
		return er
	})
	eg.Go(func() error {
		var er error
		intr, er = h.intrSvc.Get(ctx, h.biz, id, uc.Uid)
		return er
	})

	// 等待结果
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("查询文章失败，系统错误",
			logger.Int64("aid", id),
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}

	//go func() {
	// 1. 如果你想摆脱原本主链路的超时控制，你就创建一个新的
	// 2. 如果你不想，你就用 ctx
	//newCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//er := h.intrSvc.IncrReadCnt(newCtx, h.biz, art.Id)
	//if er != nil {
	//	h.l.Error("更新阅读数失败",
	//		logger.Int64("aid", art.Id),
	//		logger.Error(err))
	//}
	//}()

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,

			Content:    art.Content,
			AuthorId:   art.Author.Id,
			AuthorName: art.Author.Name,
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,

			Status: art.Status.ToUint8(),
			Ctime:  art.Ctime.Format(time.DateTime),
			Utime:  art.Utime.Format(time.DateTime),
		},
	})
}

func (h *ArticleHandler) Like(c *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
		// true 是点赞，false 是不点赞
		Like bool `json:"like"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	uc := c.MustGet("user").(jwt.UserClaims)
	var err error
	if req.Like {
		// 点赞
		err = h.intrSvc.Like(c, h.biz, req.Id, uc.Uid)
	} else {
		// 取消点赞
		err = h.intrSvc.CancelLike(c, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5, Msg: "系统错误",
		})
		h.l.Error("点赞/取消点赞失败",
			logger.Error(err),
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id))
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id  int64 `json:"id"`
		Cid int64 `json:"cid"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)

	err := h.intrSvc.Collect(ctx, h.biz, req.Id, req.Cid, uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5, Msg: "系统错误",
		})
		h.l.Error("收藏失败",
			logger.Error(err),
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) TopLikes(ctx *gin.Context) {
	interactives, err := h.intrSvc.TopLikes(ctx, "article")
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5, Msg: "系统错误",
		})
		h.l.Error("获取点赞最多文章失败",
			logger.Error(err),
		)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: slice.Map[domain.Interactive, InteractiveVo](interactives, func(idx int, src domain.Interactive) InteractiveVo {
			return InteractiveVo{
				BizId:   src.BizId,
				Biz:     src.Biz,
				LikeCnt: src.LikeCnt,
			}
		}),
	})
}
