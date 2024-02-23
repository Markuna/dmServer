package service

import (
	"douyinApi/dao"
	r "douyinApi/redis"
	"encoding/json"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type LiveComment struct {
	MsgId     string `json:"msg_id"`     // string类型id
	SecOpenId string `json:"sec_openid"` // 评论用户的加密openid, 当前其实没有加密
	Content   string `json:"content"`    // 评论内容
	AvatarUrl string `json:"avatar_url"` //评论用户头像
	Nickname  string `json:"nickname"`   // 评论用户昵称(不加密)
	Timestamp int64  `json:"timestamp"`  // 评论毫秒级时间戳
}

type LiveCommentWithRank struct {
	LiveComment
	WinCount  int64 `json:"win_count"`
	WorldRank int64 `json:"world_rank"`
}

type LiveGift struct {
	MsgId     string `json:"msg_id"`      // string类型id
	SecOpenId string `json:"sec_openid"`  // 评论用户的加密openid, 当前其实没有加密
	SecGiftId string `json:"sec_gift_id"` // 加密的礼物id
	GiftNum   int32  `json:"gift_num"`    // 送出的礼物数量
	GiftValue int32  `json:"gift_value"`  // 礼物总价值，单位分
	AvatarUrl string `json:"avatar_url"`  // 用户头像
	Nickname  string `json:"nickname"`    // 用户昵称(不加密)
	Timestamp int64  `json:"timestamp"`   // 时间戳
}

type LiveGiftWithRank struct {
	LiveGift
	WinCount  int64 `json:"win_count"`
	WorldRank int64 `json:"world_rank"`
}

type LiveLike struct {
	MsgId     string `json:"msg_id"`     // string类型id
	SecOpenId string `json:"sec_openid"` // 评论用户的加密openid, 当前其实没有加密
	LikeNum   int64  `json:"like_num"`   // 点赞数量，上游2s合并一次数据
	AvatarUrl string `json:"avatar_url"` // 用户头像
	Nickname  string `json:"nickname"`   // 用户昵称(不加密)
	Timestamp int64  `json:"timestamp"`  // 时间戳
}

func HandleMsg(msgType, roomId, bodyString string) {
	switch msgType {
	case "live_comment":
		handleLiveComment(roomId, bodyString)
	case "live_gift":
		handleLiveGift(roomId, bodyString)
	case "live_like":
		handleLiveLike(roomId, bodyString)
	}
}

func handleLiveComment(roomId, bodyString string) {
	var liveComment []LiveComment
	err := json.Unmarshal([]byte(bodyString), &liveComment)
	if err != nil {
		log.Println(err)
	}
	for _, v := range liveComment {
		if v, ok := handleWinCountAndWordRank_Comment(v); ok {
			s, _ := json.Marshal(v)
			Send(makeQueueName(roomId, "live_comment"), string(s))
			continue
		}
		s, _ := json.Marshal(v)
		Send(makeQueueName(roomId, "live_comment"), string(s))
	}
}

func handleLiveGift(roomId, bodyString string) {
	var liveGift []LiveGift
	err := json.Unmarshal([]byte(bodyString), &liveGift)
	if err != nil {
		log.Println(err)
	}
	go giftSave2Db(roomId, liveGift)
	for _, v := range liveGift {
		if v, ok := handleWinCountAndWordRank_Gift(v); ok {
			s, _ := json.Marshal(v)
			Send(makeQueueName(roomId, "live_gift"), string(s))
			continue
		}
		s, _ := json.Marshal(v)
		Send(makeQueueName(roomId, "live_gift"), string(s))
	}
}

func handleLiveLike(roomId, bodyString string) {
	var liveLike []LiveLike
	err := json.Unmarshal([]byte(bodyString), &liveLike)
	if err != nil {
		log.Println(err)
	}
	for _, v := range liveLike {
		s, _ := json.Marshal(v)
		Send(makeQueueName(roomId, "live_like"), string(s))
	}
}

func makeQueueName(roomId, msgType string) string {
	return "douyin." + msgType + "." + roomId
}

func giftSave2Db(roomId string, liveGift []LiveGift) {
	// get roomInfo from roomId

	var giftRechargeLog []dao.JkRechargeLog
	for _, v := range liveGift {
		val := decimal.NewFromFloat((float64(v.GiftValue) / 100))
		// 将 Unix 时间戳转换为 time.Time 类型
		createDateTime := time.Unix(v.Timestamp/1000, (v.Timestamp%1000)*int64(time.Millisecond))
		createDateTimeStr := createDateTime.Format("2006-01-02 15:04:05")
		giftRechargeLog = append(giftRechargeLog, dao.JkRechargeLog{
			StreamerId:   "",
			RoomId:       roomId,
			MsgId:        v.MsgId,
			PlayerId:     v.SecOpenId,
			GiftId:       v.SecGiftId,
			GiftNum:      v.GiftNum,
			GiftValue:    val,
			PlatformType: "douyin",
			CreateBy:     "system",
			CreateTime:   createDateTimeStr,
			UpdateBy:     "",
			UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
			Remark:       "",
		})
	}
	dao.InsertGiftRechargeLog(&giftRechargeLog)
}

func handleWinCountAndWordRank_Comment(liveComment LiveComment) (LiveCommentWithRank, bool) {
	joinCmds := cacheJoinCmd()
	if !contains(joinCmds, liveComment.Content) {
		return LiveCommentWithRank{}, false
	}
	rsw, _ := r.RedisDb.HGet(liveComment.SecOpenId, "winCount").Int64()

	wr, err := r.RedisDb.ZRevRank("word_rank", liveComment.SecOpenId).Result()
	if err != nil {
		return LiveCommentWithRank{
			LiveComment: liveComment,
			WinCount:    rsw,
			WorldRank:   0,
		}, true
	}
	return LiveCommentWithRank{
		LiveComment: liveComment,
		WinCount:    rsw,
		WorldRank:   wr + 1,
	}, true
}
func handleWinCountAndWordRank_Gift(liveGift LiveGift) (LiveGiftWithRank, bool) {
	rsw, _ := r.RedisDb.HGet(liveGift.SecOpenId, "winCount").Int64()
	wr, err := r.RedisDb.ZRevRank("word_rank", liveGift.SecOpenId).Result()
	if err != nil {
		return LiveGiftWithRank{
			LiveGift:  liveGift,
			WinCount:  rsw,
			WorldRank: 0,
		}, true
	}
	return LiveGiftWithRank{
		LiveGift:  liveGift,
		WinCount:  rsw,
		WorldRank: wr + 1,
	}, true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

const cacheJoinCmdKey = "douyin:joinCmd"
const expire_time = 1 * time.Minute

func cacheJoinCmd() []string {
	value := r.RedisDb.Get(cacheJoinCmdKey).Val()
	if value == "" {
		dicts := dao.QueryDictByKey("joinCmd")
		joinCmds := []string{}
		for _, v := range *dicts {
			joinCmds = append(joinCmds, v.DValue)
		}
		jsonArrs := jsonArr{Strings: joinCmds}
		newValue, err := json.Marshal(jsonArrs)
		if err != nil {
			log.Println(err)
			return []string{}
		} else {
			r.RedisDb.Set(cacheJoinCmdKey, newValue, expire_time).Err()
		}
		return joinCmds
	}
	var jsonArrs jsonArr
	err := json.Unmarshal([]byte(value), &jsonArrs)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return jsonArrs.Strings
}

type jsonArr struct {
	Strings []string `json:"strings"`
}
