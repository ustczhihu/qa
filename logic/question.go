package logic

import (
	"github.com/go-redis/redis"
	"qa/dao"
	"qa/model"
	"qa/util"
	"strconv"
)

var ZSetKey = "question"
var UpdateQuestionViewCountChan = make(chan uint64, 100)
var CreateQuestionViewCountChan = make(chan uint64, 100)
var InitQuestionViewCountChan = make(chan int, 1)

func QuestionViewCount() {
	for {
		select {
		case <-InitQuestionViewCountChan:
			qlist, code := model.GetAllQuestionId()
			if code != util.CodeSuccess {
				util.Log.Error("redis中问题的viewcount初始化失败")
				return
			}
			for _, q := range qlist {
				//score, _ := strconv.ParseFloat(strconv.Itoa(q.ViewCount), 64)
				dao.RDB.ZAdd(ZSetKey, redis.Z{Score: (float64)(q.ViewCount), Member: strconv.FormatUint(q.ID, 10)})
			}

		case qid := <-UpdateQuestionViewCountChan:
			var q model.Question
			q.ID = qid
			err := q.IncrView()
			if err != nil {
				util.Log.Error(err)
				return
			}

			err = dao.RDB.ZIncrBy(ZSetKey, 1, strconv.FormatUint(q.ID, 10)).Err()
			if err != nil {
				util.Log.Error(err)
				return
			}

		case qid := <-CreateQuestionViewCountChan:
			err := dao.RDB.ZAdd(ZSetKey, redis.Z{Score: 0, Member: strconv.FormatUint(qid, 10)}).Err()

			if err != nil {
				util.Log.Error(err)
			}
		}
	}
}
