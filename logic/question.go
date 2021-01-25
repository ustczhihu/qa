package logic

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"qa/dao"
	"qa/model"
	"qa/util"
	"strconv"
)

//热榜计算规则：score=(create_time%1e8)/1e3+view_count+5*answer_count
//score=create_time+（view_count+10*answer_count）*432
const (
	ZSetKey      = "question:zset:"
	HSetKey      = "question:hset:"
	SetKey       = "question:set:"
	HViewCount   = "view_count"
	HAnwserCount = "answer_count"
)

var InitQuestionChan = make(chan interface{})

var UpdateQuestionViewCountChan = make(chan uint64, 100)
var CreateQuestionViewCountChan = make(chan uint64, 100)
var ViewCountWrite2SQLchan = make(chan map[uint64]int)

var UpdateQuestionAnswerCountChan = make(chan uint64, 100)
var CreateQuestionAnswerCountChan = make(chan uint64, 100)
var AnswerCountWrite2SQLchan = make(chan map[uint64]int)

func QusetionRedis2Mysql() {
	for {
		select {
		//从数据库中读入(qid,view_count)/(qid,answer_count)列表
		case <-InitQuestionChan:
			qlist, code := model.GetAllQuestionId()
			if code != util.CodeSuccess {
				util.Log.Error("redis中问题的view_count初始化失败")
				return
			}
			for _, q := range qlist {
				//score, _ := strconv.ParseFloat(strconv.Itoa(q.ViewCount), 64)
				dao.RDB.HSet(HSetKey+strconv.FormatUint(q.ID, 10), HViewCount, q.ViewCount)
				dao.RDB.HSet(HSetKey+strconv.FormatUint(q.ID, 10), HAnwserCount, q.AnswerCount)
				//初始化热榜
				score := (float64)(util.Strtime2Int(q.CreatedAt))
				//fmt.Println(util.Strtime2Int(q.CreatedAt))
				dao.RDB.ZAdd(ZSetKey, redis.Z{Score: score + (float64)((q.ViewCount+10*q.AnswerCount)*432), Member: strconv.FormatUint(q.ID, 10)})
			}

		//将(qid,incrCount)列表更新到mysql中
		case set := <-ViewCountWrite2SQLchan:
			for k, v := range set {
				var q model.Question
				q.ID = k
				err := q.IncrView(v)
				if err != nil {
					util.Log.Error(err)
					return
				}
			}

		//将(qid,incrCount)列表更新到mysql中
		case set := <-AnswerCountWrite2SQLchan:
			for k, v := range set {
				var q model.Question
				q.ID = k
				err := q.IncrAnswer(v)
				if err != nil {
					util.Log.Error(err)
					return
				}
			}

		}
	}
}

func QuestionViewCount() {
	set := make(map[uint64]int)
	count := 0
	for {
		select {
		//增加浏览量
		case qid := <-UpdateQuestionViewCountChan:
			count++
			//err := dao.RDB.ZIncrBy(ZSetKey, 1, strconv.FormatUint(qid, 10)).Err()
			//e.g. key=question3109701881430017
			err := dao.RDB.HIncrBy(HSetKey+strconv.FormatUint(qid, 10), HViewCount, 1).Err()
			if err != nil {
				util.Log.Error(err)
				return
			}

			//更新热榜分数
			dao.RDB.ZIncrBy(ZSetKey, 432, strconv.FormatUint(qid, 10))

			//记录要更新的(qid,incrCount)，延迟更新
			if v, ok := set[qid]; ok {
				set[qid] = v + 1
			} else {
				set[qid] = 1
			}

			//通过channel传递要更新的(qid,incrCount)列表
			if count > 2 {
				ViewCountWrite2SQLchan <- set
				set = make(map[uint64]int)
				count = 0
			}

		//初始化新问题的浏览量
		case qid := <-CreateQuestionViewCountChan:
			//err := dao.RDB.ZAdd(ZSetKey, redis.Z{Score: 0, Member: strconv.FormatUint(qid, 10)}).Err()
			err := dao.RDB.HSet(HSetKey+strconv.FormatUint(qid, 10), HViewCount, 0).Err()
			if err != nil {
				util.Log.Error(err)
			}
		}
	}
}

func QuestionAnswerCount() {
	set := make(map[uint64]int)
	count := 0
	for {
		select {
		//增加回答量
		case qid := <-UpdateQuestionAnswerCountChan:
			count++
			//e.g. key=question3109701881430017
			err := dao.RDB.HIncrBy(HSetKey+strconv.FormatUint(qid, 10), HAnwserCount, 1).Err()
			if err != nil {
				util.Log.Error(err)
				return
			}

			//更新热榜分数
			dao.RDB.ZIncrBy(ZSetKey, 10*432, strconv.FormatUint(qid, 10))

			//记录要更新的(qid,incrCount)，延迟更新
			if v, ok := set[qid]; ok {
				set[qid] = v + 1
			} else {
				set[qid] = 1
			}

			//通过channel传递要更新的(qid,incrCount)列表
			if count > 2 {
				AnswerCountWrite2SQLchan <- set
				set = make(map[uint64]int)
				count = 0
			}

		//初始化新问题的浏览量
		case qid := <-CreateQuestionAnswerCountChan:
			err := dao.RDB.HSet(HSetKey+strconv.FormatUint(qid, 10), HAnwserCount, 0).Err()
			if err != nil {
				util.Log.Error(err)
			}
		}
	}
}

func QuestionHotList()[]model.Question{
	// 取分数最高的5个
	ret, err := dao.RDB.ZRevRangeWithScores(ZSetKey, 0, 4).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return nil
	}

	var questionList []model.Question
	for _, z := range ret {
		qid:=z.Member.(string)
		val, err := dao.RDB.Get(SetKey + qid).Result()
		var question model.Question
		if err == redis.Nil {
			qidInt,_:=strconv.Atoi(qid)
			qidInt64 := (uint64)(qidInt)
			question.ID=qidInt64
			question.Get()
			questionList=append(questionList, question)

			val,_:=json.Marshal(&question)
			dao.RDB.Set(SetKey+qid,val,0)
		}else{
			json.Unmarshal([]byte(val),&question)
			questionList=append(questionList, question)
		}
	}
	return questionList
}