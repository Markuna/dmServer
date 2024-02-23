package service

import (
	"douyinApi/dao"
	"log"

	"github.com/shopspring/decimal"
)

type StartInfo struct {
	RoomId           string `json:"roomId"`
	SuccessStartGame bool   `json:"successStartGame"`
	SuccessGetToken  bool   `json:"successGetToken"`
	SuccessGetRoomId bool   `json:"successGetRoomId"`
	SuccessGetSocket bool   `json:"successGetSocket"`
}

func (s *StartInfo) Save() {

	d := dao.StartInfo{
		RoomId:           s.RoomId,
		SuccessStartGame: convertBool(s.SuccessStartGame),
		SuccessGetToken:  convertBool(s.SuccessGetToken),
		SuccessGetRoomId: convertBool(s.SuccessGetRoomId),
		SuccessGetSocket: convertBool(s.SuccessGetSocket),
	}
	d.InsertLog()
}

func convertBool(b bool) int64 {
	if b {
		return 1
	} else {
		return 0
	}
}

type GameInfo struct {
	RoomId          string `json:"roomId"`
	GameStartTime   string `json:"gameStartTime"`
	GameEndTime     string `json:"gameEndTime"`
	PlayerNum       int64  `json:"playerNum"`
	PayUserNum      int64  `json:"payUserNum"`
	ReviveUserNum   int64  `json:"reviveUserNum"`
	DefeatChestNum  int64  `json:"defeatChestNum"`
	Gift1Num        int64  `json:"gift1Num"`
	Gift2Num        int64  `json:"gift2Num"`
	Gift3Num        int64  `json:"gift3Num"`
	Gift4Num        int64  `json:"gift4Num"`
	Gift5Num        int64  `json:"gift5Num"`
	Gift6Num        int64  `json:"gift6Num"`
	Earnings        string `json:"earnings"`
	DifficultyLevel int64  `json:"difficultyLevel"`
	LikeNum         int64  `json:"likeNum"`
}

func (s *GameInfo) Save() {
	val, err := decimal.NewFromString(s.Earnings)
	if err != nil {
		log.Println(err)
	}
	d := dao.GameInfo{
		RoomId:          s.RoomId,
		GameStartTime:   s.GameStartTime,
		GameEndTime:     s.GameEndTime,
		PlayerNum:       s.PlayerNum,
		PayUserNum:      s.PayUserNum,
		ReviveUserNum:   s.ReviveUserNum,
		DefeatChestNum:  s.DefeatChestNum,
		Gift1Num:        s.Gift1Num,
		Gift2Num:        s.Gift2Num,
		Gift3Num:        s.Gift3Num,
		Gift4Num:        s.Gift4Num,
		Gift5Num:        s.Gift5Num,
		Gift6Num:        s.Gift6Num,
		Earnings:        val,
		DifficultyLevel: s.DifficultyLevel,
		LikeNum:         s.LikeNum,
	}
	d.InsertLog()
}

type PlayerInfo struct {
	RoomId    string   `json:"roomId"`
	PlayerIds []string `json:"playerIds"`
	PlayTime  string   `json:"playTime"`
}

func (s *PlayerInfo) Save() {
	for _, v := range s.PlayerIds {
		d := dao.PlayerInfo{
			RoomId:   s.RoomId,
			PlayerId: v,
			PlayTime: s.PlayTime,
		}
		d.InsertLog()
	}
}
