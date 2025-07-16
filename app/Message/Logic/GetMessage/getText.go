package GetMessage

import (
	"common/chain"
	Redis2 "common/redis"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"message/DAO/Redis"
	"message/Models"
	"rpc/message"
	"strconv"
)

const (
	OnceMaxMessageNumber = 10
)

func RefreshText(db *gorm.DB, sender, receiver string) []*message.Text {
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), redisHandler(), mysqlHandler())
	ctx.Set("db", db)
	ctx.Set("sender", sender)
	ctx.Set("receiver", receiver)

	err := ctx.Apply()
	if err != nil {
		return nil
	}

	texts, err := chain.GetToType[[]*message.Text](ctx, "texts")
	if err != nil {
		return nil
	}

	return texts
	//	var (
	//		page     int
	//		pageSize int
	//		mText    []*message.Text
	//	)
	//	page = 1
	//	pageSize = 10
	//	// 规范化 key
	//	user1, user2 := sender, receiver
	//	u1, _ := strconv.Atoi(user1)
	//	u2, _ := strconv.Atoi(user2)
	//	if u1 > u2 {
	//		user1, user2 = user2, user1
	//	}
	//	key := fmt.Sprintf("chat:%s:%s", user1, user2)
	//
	//	// 计算分页范围
	//	start := (page - 1) * pageSize
	//	end := start + pageSize - 1
	//
	//	// 获取原始数据
	//	rawList, err := redis.ByteSlices(Redis.RedisDo(Redis2.LRANGE, key, start, end))
	//	if len(rawList) == 0 {
	//		goto MYSQL
	//	}
	//	if err != nil {
	//		log.Println(err)
	//		return nil
	//	}
	//
	//	Redis.RedisDo(Redis2.EXPIRE, key, 86400)
	//
	//	// 反序列化
	//	mText = make([]*message.Text, len(rawList))
	//	for i, raw := range rawList {
	//		mText[i] = new(message.Text)
	//		if err := json.Unmarshal(raw, mText[i]); err != nil {
	//			log.Println(fmt.Errorf("decode error at index %d: %v", i, err))
	//			return nil
	//		}
	//	}
	//	return mText
	//MYSQL:
	//	Messages := Models.GetMessages[*Models.Text](db, sender, receiver, OnceMaxMessageNumber)
	//	texts := make([]*message.Text, 0, len(Messages))
	//	for _, _message := range Messages {
	//		var text = message.Text{}
	//		text.Content = string(_message.GetContent())
	//		text.SenderAccountNum = _message.SenderAccountNum
	//		text.ReceiverAccountNum = _message.ReceiverAccountNum
	//		text.Time = _message.CreatedAt.Format("2006-01-02 15:04:05")
	//		texts = append(texts, &text)
	//	}
	//
	//	return texts
}

func redisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		var (
			page     = 1
			pageSize = 10
		)
		sender, err := chain.GetToType[string](ctx, "sender")
		if err != nil {
			return err
		}
		receiver, err := chain.GetToType[string](ctx, "receiver")
		if err != nil {
			return err
		}

		// 规范化 key
		user1, user2 := sender, receiver
		u1, _ := strconv.Atoi(user1)
		u2, _ := strconv.Atoi(user2)
		if u1 > u2 {
			user1, user2 = user2, user1
		}
		key := fmt.Sprintf("chat:%s:%s", user1, user2)

		// 计算分页范围
		start := (page - 1) * pageSize
		end := start + pageSize - 1

		// 获取原始数据
		rawList, err := redis.ByteSlices(Redis.RedisDo(Redis2.LRANGE, key, start, end))
		if len(rawList) == 0 {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = Redis.RedisDo(Redis2.EXPIRE, key, 86400)
		if err != nil {
			return err
		}

		// 反序列化
		texts := make([]*message.Text, len(rawList))
		for i, raw := range rawList {
			texts[i] = new(message.Text)
			if err := json.Unmarshal(raw, texts[i]); err != nil {
				return fmt.Errorf("decode error at index %d: %v", i, err)
			}
		}
		ctx.Set("texts", texts)
		ctx.Abort()
		return nil
	}
}

func mysqlHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		sender, err := chain.GetToType[string](ctx, "sender")
		if err != nil {
			return err
		}
		receiver, err := chain.GetToType[string](ctx, "receiver")
		if err != nil {
			return err
		}
		db, err := chain.GetToType[*gorm.DB](ctx, "db")
		if err != nil {
			return err
		}

		Messages := Models.GetMessages[*Models.Text](db, sender, receiver, OnceMaxMessageNumber)
		texts := make([]*message.Text, 0, len(Messages))
		for _, _message := range Messages {
			var text = message.Text{}
			text.Content = string(_message.GetContent())
			text.SenderAccountNum = _message.SenderAccountNum
			text.ReceiverAccountNum = _message.ReceiverAccountNum
			text.Time = _message.CreatedAt.Format("2006-01-02 15:04:05")
			texts = append(texts, &text)
		}
		ctx.Set("texts", texts)
		return nil
	}
}
