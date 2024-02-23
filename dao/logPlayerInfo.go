package dao

type PlayerInfo struct {
	RoomId   string `gorm:"column:room_id"`
	PlayerId string `gorm:"column:player_id"`
	PlayTime string `gorm:"column:play_time"`
}

func (PlayerInfo) TableName() string {
	return "dm_log_player_info"
}

func (r *PlayerInfo) InsertLog() {
	DB.Create(r)
}
