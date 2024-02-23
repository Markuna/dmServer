package dao

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

type JkRechargeLog struct {
	Id           int             `gorm:"column:id"`
	RoomId       string          `gorm:"column:room_id"`
	MsgId        string          `gorm:"column:msg_id"`
	StreamerId   string          `gorm:"column:streamer_id"`
	PlayerId     string          `gorm:"column:player_id"`
	GiftId       string          `gorm:"column:gift_id"`
	GiftNum      int32           `gorm:"column:gift_num"`
	GiftValue    decimal.Decimal `gorm:"column:gift_value;type:decimal(11,2)"`
	PlatformType string          `gorm:"column:platform_type"`
	CreateBy     string          `gorm:"column:create_by"`
	CreateTime   string          `gorm:"column:create_time"`
	UpdateBy     string          `gorm:"column:update_by"`
	UpdateTime   string          `gorm:"column:update_time"`
	Remark       string          `gorm:"column:remark"`
}

func (JkRechargeLog) TableName() string {
	return "dm_recharge_log"
}

func InsertGiftRechargeLog(r *[]JkRechargeLog) {
	DB.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(r, len(*r))
}
