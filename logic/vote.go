package logic

import (
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"qa/dao"
	"qa/model"
	"qa/util"
	"strconv"
)

/*投票的几种情况：
direction=1时，有两种情况：
	1.之前没有投过票，现在投赞成票
	2.之前投反对票，现在改成赞成票
direction=0时，两种情况：
	1.之前投过赞成，现在取消
	2.之前投过反对，现在取消
direction=-1时，两种情况：
	1.之前没有投过票，现在投反对票
	2.之前投赞成票，现在投反对票
*/

const (
	//存储某个answer有哪些用户点踩
	VoterListZSetKey   = "vote:zset:voterlist:" //+answer_id
	//存储某个question下answer的排名，根据score
	AnswerScoreZSetKey = "vote:zset:answerscore:" //+question_id
)

var InitVoteChan = make(chan interface{})
var UpdateVoteChan = make(chan model.AnswerVote, 100)
var CreateVoteChan = make(chan model.Answer, 100)
var VoteWrite2SQLchan = make(chan model.AnswerVote, 100)

func VoteForAnswer() {
	for {
		select {
		case answer:=<-CreateVoteChan:
			qid:=strconv.FormatUint(answer.QuestionID,10)
			aid:=strconv.FormatUint(answer.ID,10)
			dao.RDB.ZAdd(AnswerScoreZSetKey+qid, redis.Z{
				Member:aid,
				Score: 0,
			})

		case answerVote := <-UpdateVoteChan:
			questionID :=strconv.FormatUint(answerVote.QuestionID,10)
			answerID := strconv.FormatUint(answerVote.AnswerID, 10)
			userID := strconv.FormatUint(answerVote.UserID, 10)
			direction := answerVote.Direction

			ov := dao.RDB.ZScore(VoterListZSetKey+answerID, userID).Val()
			var op float64
			if direction > ov {
				op = 1
			} else {
				op = -1
			}
			diff := math.Abs(ov - direction)

			pipeline := dao.RDB.TxPipeline()
			//更新分数
			pipeline.ZIncrBy(AnswerScoreZSetKey+questionID, op*diff*scorePerView, answerID).Result()

			if direction == 0 {
				pipeline.ZRem(VoterListZSetKey+answerID, userID).Result()
			} else {
				pipeline.ZAdd(VoterListZSetKey+answerID, redis.Z{
					Score:  direction,
					Member: userID,
				}).Result()
			}
			_, err := pipeline.Exec()
			if err != nil {
				util.Log.Error(err)
				return
			}
			//通知mysql更新
			VoteWrite2SQLchan<-answerVote
		}
	}
}

func VoteRedis2Mysql() {
	for {
		select {
		case <-InitVoteChan:
			alist,_:=model.GetAllAnswerId()
			for _, answer := range alist {
				//初始化
				qid:=strconv.FormatUint(answer.QuestionID,10)
				aid:=strconv.FormatUint(answer.ID,10)
				dao.RDB.ZAdd(AnswerScoreZSetKey+qid, redis.Z{
					Member:aid,
					Score: 0,
				})
			}

			vlist, code := model.GetAllVoteId()
			if code != util.CodeSuccess {
				util.Log.Error("redis中votelist初始化失败")
				return
			}
			for _, answerVote := range vlist {
				//初始化
				questionID :=strconv.FormatUint(answerVote.QuestionID,10)
				answerID := strconv.FormatUint(answerVote.AnswerID, 10)
				userID := strconv.FormatUint(answerVote.UserID, 10)
				direction := answerVote.Direction

				pipeline := dao.RDB.TxPipeline()
				pipeline.ZAdd(VoterListZSetKey+answerID, redis.Z{
					Score:  direction,
					Member: userID,
				}).Result()

				pipeline.ZIncrBy(AnswerScoreZSetKey+questionID, direction*scorePerView, answerID).Result()
				_, err := pipeline.Exec()

				if err != nil {
					util.Log.Error(err)
					return
				}
			}

		case answerVote:= <- VoteWrite2SQLchan:
			err:=answerVote.Update()
			if err != nil {
				util.Log.Error(err)
				//fmt.Println("!!!!!!!!!!!!!!!",err)
				return
			}
		}
	}
}

func GetAnswerListByScore(id string,pagenum int64,pagesize int64)([]model.Answer){
	ret, err := dao.RDB.ZRevRangeWithScores(AnswerScoreZSetKey+id, (pagenum-1)*pagesize, (pagenum-1)*pagesize+pagesize).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return nil
	}

	var answerList []model.Answer
	for _, z := range ret {
		aid:=z.Member.(string)
		var answer model.Answer
		aidInt,_:=strconv.Atoi(aid)
		aidInt64 := (uint64)(aidInt)
		answer.ID=aidInt64
		answer.Get()
		answerList=append(answerList, answer)
	}
	return answerList
}

func GetBestAnswer(id string)(model.Answer){
	var answer model.Answer
	ret, err := dao.RDB.ZRevRangeWithScores(AnswerScoreZSetKey+id, 0, 0).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return answer
	}
	for _, z := range ret {
		aid:=z.Member.(string)
		aidInt,_:=strconv.Atoi(aid)
		aidInt64 := (uint64)(aidInt)
		answer.ID=aidInt64
		answer.Get()
		break
	}
	return answer
}

func GetVoteInfo(aid string,useID uint64)(float64){
	op := redis.ZRangeBy{
		Min: "-1",
		Max: "1",
	}
	ret, err := dao.RDB.ZRangeByScoreWithScores(VoterListZSetKey+aid, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return 0
	}
	for _, z := range ret {
		if z.Member==strconv.FormatUint(useID,10) {
			return z.Score
		}
	}
	return 0
}
