package dao

type StartInfo struct {
	RoomId           string `gorm:"column:room_id"`
	SuccessStartGame int64  `gorm:"column:success_start_game"`
	SuccessGetToken  int64  `gorm:"column:success_get_token"`
	SuccessGetRoomId int64  `gorm:"column:success_get_room_id"`
	SuccessGetSocket int64  `gorm:"column:success_get_socket"`
}

func (StartInfo) TableName() string {
	return "dm_log_start_info"
}

func (r *StartInfo) InsertLog() {
	DB.Create(r)
}
