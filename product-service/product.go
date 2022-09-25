package boilerplategolang

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
)

type DefaultResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ProductRequest struct {
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type ProductData struct {
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	TOPIC_NAME = "TOPIC_PRODUCT"
)

func (c *ConfigAPI) ApplyRoute() *mux.Router {
	s := c.Router.PathPrefix("/api/v1/products").Subrouter()
	s.HandleFunc("/", c.Create).Methods("POST")
	return s
}

func (c *ConfigAPI) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp DefaultResponse
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		WriteErrorReponse(w, "Internal Server Error")
		return
	}

	var req ProductRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		WriteErrorReponse(w, "Internal Server Error")
		return
	}

	var product ProductData
	product.Name = req.Name
	product.Amount = req.Amount
	product.Price = req.Price
	product.Timestamp = time.Now()

	result, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		WriteErrorReponse(w, "Internal Server Error")
		return
	}

	var topic pubsub.Topic
	c.CreateTopicAndPublish(TOPIC_NAME, string(result), &topic)
	c.CreateSubscription(c.Pubsub, "TOPICPRODUCTSUBSData", &topic)

	resp.Data = nil
	resp.Message = fmt.Sprintf("Successfully create data  %s with data %s", TOPIC_NAME, result)
	resp.Success = true

	json.NewEncoder(w).Encode(resp)

}

func WriteErrorReponse(w http.ResponseWriter, message string) error {
	var resp DefaultResponse
	resp.Data = nil
	resp.Message = message
	resp.Success = false
	return json.NewEncoder(w).Encode(resp)
}
