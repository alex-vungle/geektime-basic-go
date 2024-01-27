package dao

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"gorm.io/gorm"
	"time"
)

const FollowRelationTableName = "follow_relations"

var (
	ErrFollowerNotFound = gorm.ErrRecordNotFound
)

type TableStoreFollowRelationDao struct {
	client *tablestore.TableStoreClient
}

func (t *TableStoreFollowRelationDao) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	request := &tablestore.SQLQueryRequest{
		// 可以替换成 select *
		// 这种写法有什么问题？有什么隐患？
		// SQL 注入的隐患
		// 在实践中，如果要利用前端的输入来拼接 SQL 语句，千万要小心 SQL 注入的问题
		// select id,follower,followee from follow_relations where follower = 1 OR 1 = 1 AND status = 2 OFFSET 0 LIMIT 10
		// select id,follower,followee from follow_relations where follower = 1; TRUNCATE users OR 1 = 1 AND status = 2 OFFSET 0 LIMIT 10
		// 用户登录 select * from xxx where username = %s AND password = %s;
		//Query: fmt.Sprintf("select id,follower,followee from %s where follower = %s AND status = %d OFFSET %d LIMIT %d",
		//	FollowRelationTableName, "1 OR 1 = 1", FollowRelationStatusActive, offset, limit)}
		Query: fmt.Sprintf("select id,follower,followee from %s where follower = %d AND status = %d OFFSET %d LIMIT %d",
			FollowRelationTableName, follower, FollowRelationStatusActive, offset, limit)}
	// SELECT * FROM xx WHERE id = ? // 是利用占位符的，然后传参数
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return nil, err
	}
	resultSet := response.ResultSet
	followRelations := make([]FollowRelation, 0, limit)
	for resultSet.HasNext() {
		row := resultSet.Next()
		followRelation := FollowRelation{}
		followRelation.ID, _ = row.GetInt64ByName("id")
		followRelation.Follower, _ = row.GetInt64ByName("follower")
		followRelation.Followee, _ = row.GetInt64ByName("followee")
		followRelations = append(followRelations, followRelation)
	}
	return followRelations, nil
}

func (t *TableStoreFollowRelationDao) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	cond := tablestore.NewCompositeColumnCondition(tablestore.LO_AND)
	// 更新条件，对标 WHERE 语句
	// 多个 Filter 是 AND 条件连在一起
	cond.AddFilter(tablestore.NewSingleColumnCondition("follower", tablestore.CT_EQUAL, follower))
	cond.AddFilter(tablestore.NewSingleColumnCondition("followee", tablestore.CT_EQUAL, followee))
	req := new(tablestore.UpdateRowChange)
	req.TableName = FollowRelationTableName
	// 我预期这一行数据是存在的
	// 不在的话，会报错
	req.SetCondition(tablestore.RowExistenceExpectation_EXPECT_EXIST)
	req.SetColumnCondition(cond)
	req.PutColumn("status", int64(status))
	_, err := t.client.UpdateRow(&tablestore.UpdateRowRequest{
		UpdateRowChange: req,
	})
	return err
}

func (t *TableStoreFollowRelationDao) CntFollower(ctx context.Context, uid int64) (int64, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("SELECT COUNT(follower) as cnt from %s where followee = %d AND status = %d",
			FollowRelationTableName, uid, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return 0, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return row.GetInt64ByName("cnt")
	}
	return 0, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("SELECT COUNT(followee) as cnt from %s where follower = %d AND status = %d",
			FollowRelationTableName, uid, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return 0, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return row.GetInt64ByName("cnt")
	}
	return 0, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) CreateFollowRelation(ctx context.Context, c FollowRelation) error {
	now := time.Now().UnixMilli()
	// UpdateRowRequest + RowExistenceExpectation_IGNORE
	// 可以实现一个 insert or update 的语义
	// 单纯的使用 update 或者 put，都不能达成这个效果
	req := new(tablestore.UpdateRowRequest)
	pk := &tablestore.PrimaryKey{}
	pk.AddPrimaryKeyColumn("follower", c.Follower)
	pk.AddPrimaryKeyColumn("followee", c.Followee)
	change := &tablestore.UpdateRowChange{
		TableName: FollowRelationTableName,
		// 有一个小的问题，这边其实可以不用 id, 直接用 follower 和 followee 构成一个主键
		// 如果要用 ID，你可以用自增主键
		PrimaryKey: pk,
	}
	change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	// 只能用 Int64，
	change.PutColumn("status", int64(c.Status))
	change.PutColumn("ctime", now)
	change.PutColumn("utime", now)
	req.UpdateRowChange = change
	_, err := t.client.UpdateRow(req)
	return err
}

func (t *TableStoreFollowRelationDao) FollowRelationDetail(ctx context.Context, follower, followee int64) (FollowRelation, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("select id,follower,followee from %s where follower = %d AND followee = %d AND status = %d",
			FollowRelationTableName, follower, followee, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return FollowRelation{}, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return t.rowToEntity(row), nil
	}
	return FollowRelation{}, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) rowToEntity(row tablestore.SQLRow) FollowRelation {
	var res FollowRelation
	res.ID, _ = row.GetInt64ByName("id")
	res.Follower, _ = row.GetInt64ByName("follower")
	res.Followee, _ = row.GetInt64ByName("followee")
	status, _ := row.GetInt64ByName("status")
	res.Status = uint8(status)
	res.Ctime, _ = row.GetInt64ByName("ctime")
	res.Utime, _ = row.GetInt64ByName("utime")
	return res
}

func NewTableStoreDao(client *tablestore.TableStoreClient) *TableStoreFollowRelationDao {
	return &TableStoreFollowRelationDao{
		client: client,
	}
}
