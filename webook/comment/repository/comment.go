package repository

import (
	"context"
	"database/sql"
	"gitee.com/geekbang/basic-go/webook/comment/domain"
	"gitee.com/geekbang/basic-go/webook/comment/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"time"
)

type CommentRepository interface {
	// FindByBiz 根据 ID 倒序查找
	// 并且会返回每个评论的三条直接回复
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]domain.Comment, error)
	// DeleteComment 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, comment domain.Comment) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	// GetCommentByIds 获取单条评论 支持批量获取
	GetCommentByIds(ctx context.Context, id []int64) ([]domain.Comment, error)
	GetMoreReplies(ctx context.Context, rid int64, id int64, limit int64) ([]domain.Comment, error)
}

type CachedCommentRepo struct {
	dao dao.CommentDAO
	l   logger.LoggerV1
}

func (c *CachedCommentRepo) GetMoreReplies(ctx context.Context, rid int64, minID int64, limit int64) ([]domain.Comment, error) {
	cs, err := c.dao.FindRepliesByRid(ctx, rid, minID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(cs))
	for _, cm := range cs {
		res = append(res, c.toDomain(cm))
	}
	return res, nil
}

func (c *CachedCommentRepo) FindByBiz(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]domain.Comment, error) {
	daoComments, err := c.dao.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(daoComments))
	ids := make([]int64, 0, len(daoComments))
	for _, d := range daoComments {
		cm := c.toDomain(d)
		// 只展示三条
		cm.Children = make([]domain.Comment, 0, 3)
		res = append(res)
		ids = append(ids, d.Id)
	}
	// 查找直接的子节点
	if ctx.Value("downgraded") == "true" {
		// 触发了降级
		return res, nil
	}
	// 只找三条
	subRes, err := c.dao.FindRepliesByPids(ctx, ids, 0, 3)
	if err != nil {
		// 这里做一个容错
		c.l.Error("查找回复失败", logger.Error(err))
		return res, nil
	}
	// 一般来说，因为批次都不大，所以双重循环和 map 比起来，性能也不会差
	for _, sr := range subRes {
		for _, r := range res {
			if r.Id == sr.PID.Int64 {
				r.Children = append(r.Children, c.toDomain(sr))
			}
		}
	}
	return res, nil
}

func (c *CachedCommentRepo) DeleteComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Delete(ctx, dao.Comment{
		Id: comment.Id,
	})
}

func (c *CachedCommentRepo) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Insert(ctx, c.toEntity(comment))
}

func (c *CachedCommentRepo) GetCommentByIds(ctx context.Context, ids []int64) ([]domain.Comment, error) {
	vals, err := c.dao.FindOneByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	comments := make([]domain.Comment, 0, len(vals))
	for _, v := range vals {
		comment := c.toDomain(v)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *CachedCommentRepo) toDomain(daoComment dao.Comment) domain.Comment {
	val := domain.Comment{
		Id: daoComment.Id,
		Commentator: domain.User{
			ID: daoComment.Uid,
		},
		Biz:     daoComment.Biz,
		BizID:   daoComment.BizID,
		Content: daoComment.Content,
		CTime:   time.UnixMilli(daoComment.Ctime),
		UTime:   time.UnixMilli(daoComment.Utime),
	}
	if daoComment.PID.Valid {
		val.ParentComment = &domain.Comment{
			Id: daoComment.PID.Int64,
		}
	}
	if daoComment.RootID.Valid {
		val.RootComment = &domain.Comment{
			Id: daoComment.RootID.Int64,
		}
	}
	return val
}

func (c *CachedCommentRepo) toEntity(domainComment domain.Comment) dao.Comment {
	daoComment := dao.Comment{
		Id:      domainComment.Id,
		Uid:     domainComment.Commentator.ID,
		Biz:     domainComment.Biz,
		BizID:   domainComment.BizID,
		Content: domainComment.Content,
	}
	if domainComment.RootComment != nil {
		daoComment.RootID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.RootComment.Id,
		}
	}
	if domainComment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.ParentComment.Id,
		}
	}
	daoComment.Ctime = time.Now().UnixMilli()
	daoComment.Utime = time.Now().UnixMilli()
	return daoComment
}

func NewCommentRepo(commentDAO dao.CommentDAO, l logger.LoggerV1) CommentRepository {
	return &CachedCommentRepo{
		dao: commentDAO,
		l:   l,
	}
}
