package middleware

import (
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/utils/config"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	delLock     sync.RWMutex
	autoDelTime time.Duration
	delmap      map[string]*model.Message
	delchan     chan string
)

func AutoDeleteInit() {
	autoDelTime = config.Conf.AutoDelete.Time
	delmap = make(map[string]*model.Message)
	delchan = make(chan string, 100)

	go DeleteMsgService()
}

func AutoDelete(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !config.Conf.AutoDelete.Enable {
			return next(c)
		}
		// 保存现场
		msg := &model.Message{}
		msg.StoredMessage = tele.StoredMessage{
			MessageID: strconv.Itoa(c.Message().ID),
			ChatID:    c.Chat().ID,
		}
		msg.Save()
		// 延迟触发
		DelayDelete(msg)
		// 运行下一个中间件
		return next(c)
	}
}

func DelayDelete(msg *model.Message, delay ...time.Duration) {
	go func() {
		delayTime := autoDelTime
		if len(delay) > 0 {
			delayTime = delay[0]
		}
		delmap[msg.StoredMessage.MessageID] = msg
		time.Sleep(delayTime)
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
			if duration > autoDelTime {
				err := bot.Delete(msg.StoredMessage)
				if err != nil {
					log.Println(err)
				}
				bot.Delete(msg.StoredMessage)
				msg.Delete()
			} else {
				// 修改的信息需要重新延迟删除
				newDelay := autoDelTime + time.Second - duration
				if newDelay < 0 || newDelay > autoDelTime {
					newDelay = autoDelTime
				}
				DelayDelete(msg, newDelay)
			}
			delLock.Unlock()
		case <-time.Tick(time.Minute):
			go func() {
				log.Println("DeleteMsgService: clean")
				delLock.RLock()
				msgs := model.GetAllMessageBefore(time.Now().Add(-autoDelTime - 10*time.Second))
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
