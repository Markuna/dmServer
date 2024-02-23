package service

import (
	r "douyinApi/redis"
	"encoding/json"
	"log"
	"sort"

	"github.com/go-redis/redis"
)

type Data interface{}

type PageInfo struct {
	PageNo     int64          `json:"pageNo"`
	PageSize   int64          `json:"pageSize"`
	TotalPages int64          `json:"totalPages"`
	Data       *[]interface{} `json:"data"`
}

type PageParam struct {
	PageNo   int64 `json:"pageNo"`
	PageSize int64 `json:"pageSize"`
}

type WordRankData struct {
	UserId    string `json:"userId"`
	Score     int64  `json:"score"`
	Result    int64  `json:"result"`   // current win or lose
	WinCount  int64  `json:"winCount"` // straight wins
	UserName  string `json:"userName"`
	AvatarUrl string `json:"avatarUrl"`
}
type ByScore []WordRankData

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type WordRankDataResp struct {
	UserId    string `json:"userId"`
	WinCount  int64  `json:"winCount"` // straight wins
	Rank      int64  `json:"rank"`
	WorldRank int64  `json:"worldRank"`
}

const WordRankDataKey = "word_rank"

func GetWordRankData(param *PageParam) *PageInfo {
	start := (param.PageNo - 1) * param.PageSize
	end := param.PageNo*param.PageSize - 1
	zrw := r.RedisDb.ZRevRangeWithScores(WordRankDataKey, start, end)
	result, err := zrw.Result()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	d := make([]interface{}, 0)
	for _, v := range result {
		rst, _ := r.RedisDb.HGet(v.Member.(string), "result").Int64()
		sw, _ := r.RedisDb.HGet(v.Member.(string), "winCount").Int64()
		d = append(d, WordRankData{
			UserId:    v.Member.(string),
			Score:     int64(v.Score),
			UserName:  r.RedisDb.HGet(v.Member.(string), "userName").Val(),
			AvatarUrl: r.RedisDb.HGet(v.Member.(string), "avatarUrl").Val(),
			Result:    rst,
			WinCount:  sw,
		})
	}

	return &PageInfo{
		PageNo:     param.PageNo,
		PageSize:   param.PageSize,
		TotalPages: r.RedisDb.ZCard(WordRankDataKey).Val(),
		Data:       &d,
	}
}

func UpdateWordRank(data []WordRankData) {
	for _, v := range data {
		// incr rank score
		r.RedisDb.ZIncr(WordRankDataKey, redis.Z{
			Score:  float64(v.Score),
			Member: v.UserId,
		})
		// set user info
		jdata, _ := json.Marshal(v)
		var jMap map[string]interface{}
		json.Unmarshal(jdata, &jMap)
		delete(jMap, "winCount")
		r.RedisDb.HMSet(v.UserId, jMap)

		// set win count
		if v.Result == 1 {
			r.RedisDb.HIncrBy(v.UserId, "winCount", 1)
		} else if v.Result == 0 {
			rsw, _ := r.RedisDb.HGet(v.UserId, "winCount").Int64()
			if rsw > 0 {
				r.RedisDb.HIncrBy(v.UserId, "winCount", -1)
			} else {
				r.RedisDb.HSet(v.UserId, "winCount", 0)
			}
		}
	}

}

func Get10UserWordRankInfo(data []WordRankData) []WordRankDataResp {
	sort.Sort(ByScore(data))
	result := make([]WordRankDataResp, 0)
	if len(data) > 10 {
		data = data[:10]
	}
	for i, v := range data {
		rsw, _ := r.RedisDb.HGet(v.UserId, "winCount").Int64()
		wr := r.RedisDb.ZRevRank(WordRankDataKey, v.UserId).Val()
		result = append(result, WordRankDataResp{
			UserId:    v.UserId,
			WinCount:  rsw,
			Rank:      int64(i) + 1,
			WorldRank: wr + 1,
		})
	}
	return result

}
