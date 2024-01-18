package web

import (
	"context"
	rewardv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/reward/v1"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"github.com/gin-gonic/gin"
	"time"
)

type RewardHandler struct {
	client rewardv1.RewardServiceClient
	//artClient articlev1.ArticleServiceClient
}

//func NewRewardHandler(client rewardv1.RewardServiceClient,
//	artClient articlev1.ArticleServiceClient) *RewardHandler {
//	return &RewardHandler{client: client, artClient: artClient}
//}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/reward")
	rg.POST("/detail",
		ginx.WrapBodyAndToken[GetRewardReq](h.GetReward))
	//rg.POST("/article",
	//	ginx.WrapBodyAndToken[GetRewardReq](h.GetReward))
}

type GetRewardReq struct {
	Rid int64
}

// GetReward 前端传过来一个超长的超时时间，例如说 10s
// 后端去轮询
// 可能引来巨大的性能问题
// 真正优雅的还是前端来轮询
// stream
func (h *RewardHandler) GetReward(
	ctx *gin.Context,
	req GetRewardReq,
	claims jwt.UserClaims) (ginx.Result, error) {

	for {
		newCtx, cancel := context.WithTimeout(ctx, time.Second)
		resp, err := h.client.GetReward(newCtx, &rewardv1.GetRewardRequest{
			Rid: req.Rid,
			Uid: claims.Id,
		})
		cancel()
		if err != nil {
			return ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			}, err
		}
		if resp.Status == 1 {
			continue
		}
		return ginx.Result{
			// 暂时也就是只需要状态
			Data: resp.Status.String(),
		}, nil
	}

}

type RewardArticleReq struct {
	Aid int64 `json:"aid"`
	Amt int64 `json:"amt"`
}
