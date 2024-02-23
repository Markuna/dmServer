package dao

import "github.com/shopspring/decimal"

type GameInfo struct {
	RoomId          string          `gorm:"column:room_id"`
	GameStartTime   string          `gorm:"column:game_start_time"`
	GameEndTime     string          `gorm:"column:game_end_time"`
	PlayerNum       int64           `gorm:"column:player_num"`
	PayUserNum      int64           `gorm:"column:pay_user_num"`
	ReviveUserNum   int64           `gorm:"column:revive_user_num"`
	DefeatChestNum  int64           `gorm:"column:defeat_chest_num"`
	Gift1Num        int64           `gorm:"column:gift1_num"`
	Gift2Num        int64           `gorm:"column:gift2_num"`
	Gift3Num        int64           `gorm:"column:gift3_num"`
	Gift4Num        int64           `gorm:"column:gift4_num"`
	Gift5Num        int64           `gorm:"column:gift5_num"`
	Gift6Num        int64           `gorm:"column:gift6_num"`
	Earnings        decimal.Decimal `gorm:"column:earnings;type:decimal(11,2)"`
	DifficultyLevel int64           `gorm:"column:difficulty_level"`
	LikeNum         int64           `gorm:"column:like_num"`
}

func (GameInfo) TableName() string {
	return "dm_log_game_info"
}

func (r *GameInfo) InsertLog() {
	DB.Create(r)
}
