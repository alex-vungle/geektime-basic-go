package web

import (
	rewardv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/reward/v1"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	client rewardv1.RewardServiceClient
}

func NewRewardHandler(client rewardv1.RewardServiceClient) *RewardHandler {
	return &RewardHandler{client: client}
}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/reward")
	rg.POST("/detail",
		ginx.WrapBodyAndClaims(h.GetReward))
}

type GetRewardReq struct {
	Rid int64
}

func (h *RewardHandler) GetReward(
	ctx *gin.Context,
	req GetRewardReq,
	claims jwt.UserClaims) (ginx.Result, error) {
	resp, err := h.client.GetReward(ctx, &rewardv1.GetRewardRequest{
		// 我这一次打赏的 ID
		Rid: req.Rid,
		// 要防止非法访问，我只能看到我打赏的记录
		// 我不能看到别人打赏记录
		Uid: claims.Uid,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		// 暂时也就是只需要状态
		Data: resp.Status.String(),
	}, nil
}
