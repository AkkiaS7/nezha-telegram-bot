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
	timerLock   sync.RWMutex
	autoDelTime time.Duration
	delmap      map[string]*model.Message
	timerMap    map[string]*time.Timer
	delchan     chan string
)

func AutoDeleteInit() {
	autoDelTime = config.Conf.AutoDelete.Time
	delmap = make(map[string]*model.Message)
	timerMap = make(map[string]*time.Timer)
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
		timerLock.Lock()
		key := msg.MessageID + "|" + strconv.FormatInt(msg.ChatID, 10)
		if msg.ID == 0 {
			delLock.Lock()
			if _, ok := delmap[key]; ok {
				msg = delmap[key]
			} else {
				msg.Save()
				delmap[key] = msg
			}
			delLock.Unlock()
		} else {
			delmap[key] = msg
		}
		delayTime := autoDelTime
		if len(delay) > 0 {
			delayTime = delay[0]
		}

		if timerMap[key] != nil {
			timerMap[key].Reset(delayTime)
			timerLock.Unlock()
			return
		}

		timer := time.AfterFunc(delayTime, func() {
			timerLock.Lock()
			defer timerLock.Unlock()
			delchan <- key
			delete(timerMap, key)
		})
		timerMap[key] = timer

		timerLock.Unlock()
	}()
}

func DeleteMsgService() {
	log.Println("DeleteMsgService start")
	for {
		select {
		case key := <-delchan:
			delLock.Lock()
			msg := delmap[key]
			delete(delmap, key)
			if msg == nil {
				delLock.Unlock()
				continue
			}
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
					cmsg := msg
					DelayDelete(&cmsg, time.Second)
				}
				delLock.RUnlock()
			}()
		}
	}
}
