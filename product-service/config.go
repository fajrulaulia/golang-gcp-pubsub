package boilerplategolang

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
)

type ConfigAPI struct {
	Router *mux.Router
	Pubsub *pubsub.Client
}

func (c *ConfigAPI) SetupGooglePubsub(ctx context.Context) *pubsub.Client {

	projectID := "for-learning-363517"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	c.Pubsub = client
	return client
}
