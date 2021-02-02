package logic

import (
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
var VoteWrite2SQLchan = make(chan model.AnswerVote, 100)

func VoteForAnswer() {
	for {
		select {
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
				return
			}
		}
	}
}
