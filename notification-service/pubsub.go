package boilerplategolang

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

func (c *ConfigAPI) CreateTopicAndPublish(topicID string, message string, topic *pubsub.Topic) {
	ctx := context.Background()

	topic = c.Pubsub.Topic(topicID)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})

	_, err := result.Get(ctx)
	if err != nil {
		log.Println("PublishMessage", err)
		os.Exit(0)
	}
}
func (c *ConfigAPI) CreateSubscription(client *pubsub.Client, name string, topic *pubsub.Topic) error {
	ctx := context.Background()
	sub, err := client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 1 * time.Minute,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	return nil
}

func (c *ConfigAPI) PullMsgs(subID string, data chan Data) error {
	ctx := context.Background()

	sub := c.Pubsub.Subscription(subID)

	var received int32
	err := sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var datas Data
		json.Unmarshal(msg.Data, &datas)
		data <- datas
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	fmt.Printf("Received %d messages\n", received)

	return nil
}
