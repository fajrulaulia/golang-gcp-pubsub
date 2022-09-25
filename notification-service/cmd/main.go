package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	app "github.com/fajrulaulia/boilerplate-golang"
	"github.com/gorilla/mux"
)

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/fajrulaulia/Workspaces/golang-google-pubsub/product-service/credentials/for-learning-363517-2f2d5c65924a.json")
	ctx := context.Background()

	r := mux.NewRouter()

	var c app.ConfigAPI
	c.Router = r
	client := c.SetupGooglePubsub(ctx)
	defer client.Close()
	c.ApplyRoute()

	if len(os.Getenv("PORT")) < 1 {
		os.Setenv("PORT", "8001")
	}
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Service on port", os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}
