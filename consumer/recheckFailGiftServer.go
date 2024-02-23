package consumer

import (
	"douyinApi/cmd/api"
	"douyinApi/config"
	"douyinApi/dao"
	"time"

	"github.com/shopspring/decimal"
)

type RoomStruct struct {
	RoomId   string
	PageSize int64 // 每页条数 固定 100条
	Offset   int64 // 偏移量（已经读取的条数）
}

// 需要被检查的RoomId
type NeedCheckRoomId struct {
	RoomSet map[string]*RoomStruct
}

// 已经断开的RoomId
type DisconnectRoomId struct {
	RoomSet map[string]interface{}
}

var NeedCheckRoomIdInstance NeedCheckRoomId
var DisconnectRoomIdInstance DisconnectRoomId

func init() {
	NeedCheckRoomIdInstance = NeedCheckRoomId{
		RoomSet: make(map[string]*RoomStruct),
	}
	DisconnectRoomIdInstance = DisconnectRoomId{
		RoomSet: make(map[string]interface{}),
	}
	// 每隔300ms拿个roomId调用下
	go func() {
		for {
			for roomId, v := range NeedCheckRoomIdInstance.RoomSet {
				// 循环调用接口，获取下一页的数据
				pageIndex := v.Offset/v.PageSize + 1
				resp := api.GetFailData(api.FailDataReq{
					RoomId:   roomId,
					AppId:    config.Get().Douyin.AppId,
					MsgType:  api.LiveGift,
					PageNum:  pageIndex,
					PageSize: v.PageSize,
				})
				if resp != nil && resp.DataList != nil {
					var giftRechargeLog []dao.JkRechargeLog
					for _, v := range resp.DataList {
						for _, v := range v.Payload {
							val := decimal.NewFromFloat((float64(v.GiftValue) / 100))
							createDateTime := time.Unix(v.Timestamp/1000, (v.Timestamp%1000)*int64(time.Millisecond))
							createDateTimeStr := createDateTime.Format("2006-01-02 15:04:05")
							giftRechargeLog = append(giftRechargeLog, dao.JkRechargeLog{
								RoomId:       roomId,
								MsgId:        v.MsgId,
								StreamerId:   "",
								PlayerId:     v.SecOpenId,
								GiftId:       v.SecGiftId,
								GiftNum:      v.GiftNum,
								GiftValue:    val,
								PlatformType: "douyin",
								CreateBy:     "failAPI",
								CreateTime:   createDateTimeStr,
								UpdateBy:     "",
								UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
								Remark:       "",
							})
						}
					}
					dao.InsertGiftRechargeLog(&giftRechargeLog)
					// 更新 offset
					v.Offset = (pageIndex-1)*v.PageSize + int64(len(giftRechargeLog))
					// 如果数据的数量和offset一致,说明已经没有数据了
					// 并且roomId在DisconnectRoomIdInstance.RoomSet中
					if v.Offset == resp.TotalCount && DisconnectRoomIdInstance.ContainsKey(roomId) {
						// 则可以从两个set里移除这个roomId
						delete(NeedCheckRoomIdInstance.RoomSet, roomId)
						delete(DisconnectRoomIdInstance.RoomSet, roomId)
					}
				}
				time.Sleep(time.Millisecond * 300)
			}
			time.Sleep(time.Millisecond * 300)
		}
	}()
}

func SetNeedCheckRoomId(roomId string) {
	if _, ok := NeedCheckRoomIdInstance.RoomSet[roomId]; ok {
		return
	} else {
		NeedCheckRoomIdInstance.RoomSet[roomId] = &RoomStruct{
			RoomId:   roomId,
			PageSize: 100,
			Offset:   0,
		}
		delete(DisconnectRoomIdInstance.RoomSet, roomId)
	}
}

func SetDisconnectRoomId(roomId string) {
	if _, ok := DisconnectRoomIdInstance.RoomSet[roomId]; ok {
		return
	} else {
		DisconnectRoomIdInstance.RoomSet[roomId] = 1
	}
}

func (d *DisconnectRoomId) ContainsKey(roomId string) bool {
	if _, ok := d.RoomSet[roomId]; ok {
		return true
	} else {
		return false
	}
}
