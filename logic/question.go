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

var InitQuestionViewCountChan = make(chan interface{})
var Write2SQLchan = make(chan map[uint64]int)

func QuestionViewCount() {
	set := make(map[uint64]int)
	for {
		select {
		//viewcount增加
		case qid := <-UpdateQuestionViewCountChan:
			err := dao.RDB.ZIncrBy(ZSetKey, 1, strconv.FormatUint(qid, 10)).Err()
			if err != nil {
				util.Log.Error(err)
				return
			}
			if v, ok := set[qid]; ok {
				set[qid] = v + 1
			} else {
				set[qid] = 1
			}

			if len(set) > 2 {
				Write2SQLchan <- set
				set = make(map[uint64]int)
			}

		//初始化新问题的viewcount
		case qid := <-CreateQuestionViewCountChan:
			err := dao.RDB.ZAdd(ZSetKey, redis.Z{Score: 0, Member: strconv.FormatUint(qid, 10)}).Err()
			if err != nil {
				util.Log.Error(err)
			}
		}
	}
}

func QusetionViewCount2Mysql() {
	for {
		select {
		//初始化viewcount列表
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

		//存储到mysql中
		case set := <-Write2SQLchan:
			for k, v := range set {
				var q model.Question
				q.ID = k
				err := q.IncrView(v)
				if err != nil {
					util.Log.Error(err)
					return
				}
			}
		}
	}
}
