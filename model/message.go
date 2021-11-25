package model

import (
	"github.com/0xJacky/Homework-api/settings"
)

type Message struct {
	Model
	FromUserID uint   `json:"from_user_id" gorm:"default:NULL"`
	From       *User  `json:"from,omitempty" gorm:"foreignKey:FromUserID;references:ID"`
	ToUserID   uint   `json:"to_user_id" gorm:"default:NULL"`
	To         *User  `json:"to,omitempty" gorm:"foreignKey:ToUserID;references:ID"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Flag       uint   `json:"flag" gorm:"default:0"`
}

// SendMessage 发送消息
func SendMessage(fromId, toId uint, title, content string) {
	// 发送
	message := Message{
		FromUserID: fromId,
		ToUserID:   toId,
		Title:      title,
		Content:    content,
		Flag:       0,
	}

	db.Create(&message)
	// 更新未读消息数量
	CountUnreadMessage(toId)
}

func BatchSendMessage(fromId uint, toId []uint, title, content string) {
	var messages []Message
	for i := range toId {
		messages = append(messages, Message{
			FromUserID: fromId,
			ToUserID:   toId[i],
			Title:      title,
			Content:    content,
			Flag:       0,
		})
	}
	db.Create(&messages)
	// 更新未读数量
	for i := range toId {
		CountUnreadMessage(toId[i])
	}
}

type MessageListItem struct {
	Model
	FromUserID   uint   `json:"from_user_id"`
	FromUsername string `json:"from_username"`
	Title        string `json:"title"`
	Flag         uint   `json:"flag"`
}

func GetMessageList(userId interface{}, pageNum int) (list []MessageListItem, total int64) {
	db.Model(&Message{}).Where("to_user_id", userId).Count(&total)

	db.Model(&Message{}).Select("messages.id, from_user_id, title, flag, messages.created_at,"+
		" messages.updated_at, users.name as from_username").
		Joins("left join users on users.id = messages.from_user_id").
		Offset(pageNum).
		Limit(settings.AppSettings.PageSize).
		Where("to_user_id", userId).
		Scan(&list)
	return
}

func ReadMessage(userId, messageId interface{}) (m Message) {

	db.Model(&Message{}).
		Where("to_user_id", userId).
		Where("id", messageId).Update("flag", 1)

	CountUnreadMessage(userId)

	db.Joins("From").First(&m, messageId)

	return
}

func DeleteMessage(userId, messageId interface{}) {

	db.Where("to_user_id", userId).Delete(&Message{}, messageId)

	CountUnreadMessage(userId)
}

func DeleteAllMessages(userId interface{}) {
	db.Delete(&Message{}, "to_user_id", userId)

	CountUnreadMessage(userId)
}

// CountUnreadMessage 更新未读消息数量
func CountUnreadMessage(userId interface{}) {
	var unreadNum int64
	db.Model(&Message{}).Where("to_user_id", userId).
		Where("flag", 0).
		Count(&unreadNum)

	db.Model(&User{}).Where("id", userId).Update("unread_message_num", unreadNum)
}
