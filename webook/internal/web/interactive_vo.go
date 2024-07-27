package web

type InteractiveVo struct {
	Biz        string `json:"biz,omitempty"`
	BizId      int64  `json:"id,omitempty"`
	ReadCnt    int64  `json:"read_count,omitempty"`
	LikeCnt    int64  `json:"like_count,omitempty"`
	CollectCnt int64  `json:"collect_count,omitempty"`
	Liked      bool   `json:"liked,omitempty"`
	Collected  bool   `json:"collected,omitempty"`
}
