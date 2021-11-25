package live

import (
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/rds"
	"github.com/spf13/cast"
	"log"
	"strings"
	"time"
)

func SyncLastActive()  {
	for  {
		log.Println("[SyncLastActive] 同步用户上次活跃时间开始")

		keys, err := rds.Keys("lastActive:user:*")
		if err != nil {
			log.Println(err)
		}
		log.Println(keys)

		var value string
		for i := range keys {
			key := keys[i]
			value, err = rds.GetNoPrefix(key)

			keySplit := strings.Split(key, ":")
			if len(keySplit) == 4 {
				userId := keySplit[3]
				log.Println(userId, value)

				user := model.NewUser(userId)
				t := cast.ToTime(value)
				user.UpdatesWithoutPreload(&model.User{
					LastActive: &t,
				})
			}
			err = rds.DelNoPrefix(key)

			if err != nil {
				log.Println(err)
			}
		}

		log.Println("[SyncLastActive] 同步用户上次活跃时间结束")

		time.Sleep(100*time.Second)
	}
}
