package controller

import (
	"encoding/json"
	"github.com/alexeysoshin/SmsBird/model"
	"github.com/alexeysoshin/SmsBird/service"
	"log"
	"math"
	"net/http"
	"strings"
)

type Controller struct {
	smsService *service.SmsService
}

// Don't allow phone numbers with less that 6 digits
var MinRecipientLength = int(math.Pow(10, 6))

const MaxOriginatorLength = 11

func NewController(key string) *Controller {
	c := &Controller{}

	c.smsService = service.NewService(key)

	return c
}

func (c *Controller) SendSms(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Println("Wrong header ")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	message := &model.Message{}
	err := json.NewDecoder(r.Body).Decode(message)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(message.Message)) == 0 {
		log.Println("Empty message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(message.Originator) == 0 || len(message.Originator) > MaxOriginatorLength {
		log.Println("Bad originator")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if message.Recipient < MinRecipientLength {
		log.Println("Bad recipient")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.smsService.SendSms(message)
}
