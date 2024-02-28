package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func (g *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).
		Select("count(follower)").
		// 我这个怎么办？
		// 考虑在 followee 上创建一个索引（followee, status)
		// <followee, follower> 行不行？
		Where("followee = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).
		Select("count(followee)").
		// 我在这里，能利用到 <follower, followee> 的联合唯一索引
		Where("follower = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).
		Where("follower = ? AND followee = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"utime":  now,
		}).Error
}

func (g *GORMFollowRelationDAO) FollowRelationList(ctx context.Context,
	follower, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	// 这样就达成了覆盖索引的效果
	err := g.db.WithContext(ctx).Select("follower, followee").
		// 这个查询要求我们要在 follower 上创建一个索引，或者 <follower, followee> 联合唯一索引
		// 进一步考虑，将 status 也加入索引
		Where("follower = ? AND status = ?", follower, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND followee = ? AND status = ?",
		follower, followee, FollowRelationStatusActive).First(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	now := time.Now().UnixMilli()
	f.Utime = now
	f.Ctime = now
	f.Status = FollowRelationStatusActive
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"utime":  now,
			"status": FollowRelationStatusActive,
		}),
	}).Create(&f).Error
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDao {
	return &GORMFollowRelationDAO{
		db: db,
	}
}
