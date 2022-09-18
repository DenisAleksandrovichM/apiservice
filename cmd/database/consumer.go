package main

import (
	"context"
	"github.com/DenisAleksandrovichM/apiservice/internal/database/config"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/database/core/user"
	modelsPkg "github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	eventCreate = "Create"
	eventUpdate = "Update"
	eventDelete = "Delete"
	duration    = time.Second * 10
)

type Consumer struct {
	user     userPkg.User
	commands map[string]func(context.Context, interface{}) (modelsPkg.User, error)
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			log.Info("context done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Error("data channel closed")
				return nil
			}
			ctx := context.Background()
			switch string(msg.Key) {
			case eventDelete:
				var kafkaMsg modelsPkg.KafkaMessage
				err := kafkaMsg.UnmarshalBinary(msg.Value)
				if err != nil {
					return err
				}
				err = c.user.Delete(ctx, string(kafkaMsg.Data))
				if err != nil {
					log.Error(err)
					return err
				}
				c.user.SetRedisData(kafkaMsg.Id, modelsPkg.User{})
			case eventCreate, eventUpdate:
				var kafkaMsg modelsPkg.KafkaMessage
				err := kafkaMsg.UnmarshalBinary(msg.Value)
				if err != nil {
					return err
				}
				var user modelsPkg.User
				err = user.UnmarshalBinary(kafkaMsg.Data)
				if err != nil {
					return err
				}
				var responseUser modelsPkg.User
				if string(msg.Key) == eventCreate {
					responseUser, err = c.user.Create(ctx, user)
				} else {
					responseUser, err = c.user.Update(ctx, user)
				}
				if err != nil {
					log.Error(err)
					return err
				}
				c.user.SetRedisData(kafkaMsg.Id, responseUser)
			}
		}
	}
}

func runConsumer(user userPkg.User, errSignals chan error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(config.Brokers, "startConsuming", cfg)
	if err != nil {
		errSignals <- errors.Wrapf(err, "on start consuming")
		return
	}
	consumer := &Consumer{user: user}
	for {
		if err = client.Consume(context.Background(), config.Topics, consumer); err != nil {
			log.Error("on consume: %v", err)
			waitForConnection()
		}
	}
}

func waitForConnection() {
	time.Sleep(duration)
}
