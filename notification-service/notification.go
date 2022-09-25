package boilerplategolang

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DefaultResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	TOPIC_NAME = "TOPIC_PRODUCT"
)

type Data struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *ConfigAPI) ApplyRoute() *mux.Router {
	c.InitService()
	return mux.NewRouter()
}

func (c *ConfigAPI) InitService() {
	log.Println("Running send message")

	data := make(chan Data, 1)

	go func() {
		for {
			c.PullMsgs("TOPICPRODUCTSUBSData", data)
		}
	}()

	go func() {
		for {

			log.Print("data", <-data)
			res := <-data
			c.Create(res.Name)
		}
	}()

}

func (c *ConfigAPI) Create(name string) {

	log.Println("Successfully create", name)

}
func WriteErrorReponse(w http.ResponseWriter, message string) error {
	var resp DefaultResponse
	resp.Data = nil
	resp.Message = message
	resp.Success = false
	return json.NewEncoder(w).Encode(resp)
}
