package SendMessage

import (
	"common/chain"
	Redis2 "common/redis"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"message/DAO/Redis"
	"message/Models"
	"rpc/message"
	"strconv"
	"time"
)

const (
	MaxMessageNumber = 100
)

func SendText(db *gorm.DB, sender string, receiver string, size int64, content []byte, CreatedAt time.Time) error {
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), textRedisHandler(), textMySQLHandler())
	ctx.Set("sender", sender)
	ctx.Set("receiver", receiver)
	ctx.Set("created_at", CreatedAt)
	ctx.Set("content", content)
	ctx.Set("size", size)
	ctx.Set("db", db)
	return ctx.Apply()

	//// 规范化 key (保证user1 < user2)
	//var mText message.Text
	//mText.SenderAccountNum = sender
	//mText.ReceiverAccountNum = receiver
	//mText.Time = CreatedAt.Format("2006-01-02 15:04:05")
	//mText.Content = string(content)
	//user1, user2 := sender, receiver
	//u1, _ := strconv.Atoi(user1)
	//u2, _ := strconv.Atoi(user2)
	//if u1 > u2 {
	//	user1, user2 = user2, user1
	//}
	//key := fmt.Sprintf("chat:%s:%s", user1, user2)
	//
	//// TODO: Redis
	//data, err := json.Marshal(&mText)
	//if err != nil {
	//	return err
	//}
	//
	//// 存储并设置过期时间（可选）
	//if _, err := Redis.RedisDo(Redis2.LPUSH, key, data); err != nil {
	//	return err
	//}
	//Redis.RedisDo(Redis2.LTRIM, key, 0, MaxMessageNumber)
	////一天过期
	//if _, err := Redis.RedisDo(Redis2.EXPIRE, key, 86400); err != nil { // 7天过期
	//	log.Printf("Set expire failed: %v", err)
	//}
	//
	//// TODO: MySQL
	//text := Models.NewText(size)
	//text.Content = append(text.Content, content...)
	//text.SenderAccountNum = sender
	//text.ReceiverAccountNum = receiver
	//text.CreatedAt = CreatedAt
	//Models.WriteMessage(db, text)
	//
	//return err
}

func textRedisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		sender, err := chain.GetToType[string](ctx, "sender")
		if err != nil {
			return err
		}
		receiver, err := chain.GetToType[string](ctx, "receiver")
		if err != nil {
			return err
		}
		CreatedAt, err := chain.GetToType[time.Time](ctx, "created_at")
		if err != nil {
			return err
		}
		content, err := chain.GetToType[[]byte](ctx, "content")
		if err != nil {
			return err
		}

		// 规范化 key (保证user1 < user2)
		var mText message.Text
		mText.SenderAccountNum = sender
		mText.ReceiverAccountNum = receiver
		mText.Time = CreatedAt.Format("2006-01-02 15:04:05")
		mText.Content = string(content)
		user1, user2 := sender, receiver
		u1, _ := strconv.Atoi(user1)
		u2, _ := strconv.Atoi(user2)
		if u1 > u2 {
			user1, user2 = user2, user1
		}
		key := fmt.Sprintf("chat:%s:%s", user1, user2)

		// TODO: Redis
		data, err := json.Marshal(&mText)
		if err != nil {
			return err
		}

		// 存储并设置过期时间（可选）
		if _, err := Redis.RedisDo(Redis2.LPUSH, key, data); err != nil {
			return errors.New(fmt.Sprintf("Redis LPUSH failed: %v", err))
		}
		if _, err = Redis.RedisDo(Redis2.LTRIM, key, 0, MaxMessageNumber); err != nil {
			return errors.New(fmt.Sprintf("Redis LTRIM failed: %v", err))
		}
		//一天过期
		if _, err := Redis.RedisDo(Redis2.EXPIRE, key, 86400); err != nil { // 7天过期
			return errors.New(fmt.Sprintf("Set expire failed: %v", err))
		}

		return nil
	}
}

func textMySQLHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		sender, err := chain.GetToType[string](ctx, "sender")
		if err != nil {
			return err
		}
		receiver, err := chain.GetToType[string](ctx, "receiver")
		if err != nil {
			return err
		}
		CreatedAt, err := chain.GetToType[time.Time](ctx, "created_at")
		if err != nil {
			return err
		}
		content, err := chain.GetToType[[]byte](ctx, "content")
		if err != nil {
			return err
		}
		size, err := chain.GetToType[int64](ctx, "size")
		if err != nil {
			return err
		}
		db, err := chain.GetToType[*gorm.DB](ctx, "db")
		if err != nil {
			return err
		}

		// TODO: MySQL
		text := Models.NewText(size)
		text.Content = append(text.Content, content...)
		text.SenderAccountNum = sender
		text.ReceiverAccountNum = receiver
		text.CreatedAt = CreatedAt
		Models.WriteMessage(db, text)
		return nil
	}
}
