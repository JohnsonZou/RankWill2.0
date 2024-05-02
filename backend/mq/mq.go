package mq

import (
	"RankWillServer/backend/util"
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func InitMQChanel(ctx context.Context) (context.Context, error) {

	viper := viper.New()
	viper.SetConfigName("rabbitmq_config")
	viper.SetConfigType("yaml")
	dir, _ := os.Getwd()
	viper.AddConfigPath(dir + "\\config\\")
	if err := viper.ReadInConfig(); err != nil {
		return ctx, err
	}
	conn, err := amqp.Dial(viper.GetString("dialstr"))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare(
		"delayedExchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct", // 延迟交换机
		},
	)
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		"delayedQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind(
		q.Name,
		"",
		"delayedExchange",
		false,
		nil,
	)

	return context.WithValue(ctx, util.MainMQChanelKey, ch), nil
}
func SendMsgToDelayQueueByDelaySeconds(ctx context.Context, msg string, t int64) error {
	ch := util.GetMainMQChanel(ctx)
	if ch == nil {
		log.Println("[MQ]nil chanel")
		return errors.New("[MQ]nil chanel")
	}

	err := ch.Publish(
		"delayedExchange", // exchange
		"",                // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
			Headers: amqp.Table{
				"x-delay": t * 1000, // 设置延迟时间
			},
		})
	log.Println("a msg :", msg, "is about to arrive at ", t, " seconds later, err: ", err)
	return err
}
func SendMsgToDelayQueueByUnixTime(ctx context.Context, msg string, t int64) error {
	tDel := t - time.Now().Unix()
	if tDel <= 0 {
		return errors.New("invalid msg time")
	}
	return SendMsgToDelayQueueByDelaySeconds(ctx, msg, tDel)
}
