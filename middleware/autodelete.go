package middleware

import (
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"sync"
	"time"
)

var delLock = sync.RWMutex{}
var delmap = make(map[string]*model.Message)
var delchan = make(chan string, 100)

func AutoDelete(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		// 保存现场
		msg := &model.Message{}
		msg.StoredMessage = tele.StoredMessage{
			MessageID: strconv.Itoa(c.Message().ID),
			ChatID:    c.Chat().ID,
		}
		msg.Save()
		// 延迟触发
		DelayDelete(msg, time.Second*1)
		// 运行下一个中间件
		return next(c)
	}
}

func DelayDelete(msg *model.Message, delay time.Duration) {
	go func() {
		delmap[msg.StoredMessage.MessageID] = msg
		time.Sleep(delay)
		log.Println("add to chan:", msg.StoredMessage.MessageID)
		delchan <- msg.StoredMessage.MessageID
	}()
}

func DeleteMsgService() {
	log.Println("DeleteMsgService start")
	for {
		select {
		case msgID := <-delchan:
			delLock.Lock()
			msg := delmap[msgID]
			delete(delmap, msgID)
			if msg == nil {
				continue
			}
			log.Println("DeleteMsgService:", msg.StoredMessage.MessageID)
			//calc sec between updateTime and now
			duration := time.Now().Sub(msg.UpdatedAt)
			if duration > 10 {
				err := bot.Delete(msg.StoredMessage)
				if err != nil {
					log.Println(err)
				}
				bot.Delete(msg.StoredMessage)
				msg.Delete()
			} else {
				// 修改的信息需要重新延迟删除
				newDelay := time.Second*11 - duration
				if newDelay < 0 || newDelay > time.Second*10 {
					newDelay = time.Second * 10
				}
				DelayDelete(msg, newDelay)
			}
			delLock.Unlock()
		case <-time.Tick(time.Minute):
			go func() {
				log.Println("DeleteMsgService: clean")
				delLock.RLock()
				msgs := model.GetAllMessageBefore(time.Now().Add(-time.Second * 30))
				for _, msg := range msgs {
					delete(delmap, msg.StoredMessage.MessageID)
					bot.Delete(msg.StoredMessage)
					msg.Delete()
				}
				delLock.RUnlock()
			}()
		}
	}
}
