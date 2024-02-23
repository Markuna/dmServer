package proto

type PushMsg struct {
	PushType string   `json:"push_type"`
	Payload  *Payload `json:"payload"`
}

type Payload struct {
	MsgId     string `json:"msg_id"`      // string类型id
	SecOpenid string `json:"sec_openid"`  // 评论用户的加密openid, 当前其实没有加密
	Content   string `json:"content"`     // 评论内容
	SecGiftId string `json:"sec_gift_id"` // 加密的礼物id
	GiftNum   int32  `json:"gift_num"`    // 送出的礼物数量
	GiftValue int32  `json:"gift_value"`  // 礼物总价值，单位分
	LikeNum   int64  `json:"like_num"`    // 点赞数量，上游2s合并一次数据
	AvatarUrl string `json:"avatar_url"`  // 评论用户头像
	Nickname  string `json:"nickname"`    // 评论用户昵称(不加密)
	Timestamp int64  `json:"timestamp"`   // 评论毫秒级时间戳

	WinCount  int32 `json:"win_count"`  // 连胜数
	WorldRank int32 `json:"world_rank"` // 世界排名
}
