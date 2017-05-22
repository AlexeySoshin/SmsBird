package service

import (
	"fmt"
	"github.com/alexeysoshin/SmsBird/model"
	"github.com/messagebird/go-rest-api"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type SmsService struct {
	timer  <-chan time.Time
	queue  chan model.Message
	done   <-chan bool
	client *messagebird.Client
}

const MaxLength = 160

func NewService(key string) *SmsService {

	s := &SmsService{}

	s.timer = time.NewTicker(time.Second * 1).C
	s.queue = make(chan model.Message, 1000)
	s.done = make(chan bool)
	s.client = messagebird.New(key)

	go startSend(s.timer, s.queue, s.done, s)

	return s
}

// Send one or more SMS messages
func (s *SmsService) SendSms(m *model.Message) {

	messages := []string{}
	if len(m.Message) > MaxLength {
		messages = splitMessage(m.Message)
	} else {
		messages = append(messages, m.Message)
	}

	// Reference should be the same for all messages
	ref := rand.Intn(255)

	for i, text := range messages {
		log.Println("Enqueuing message")
		s.queue <- model.Message{
			Originator: m.Originator,
			Message:    text,
			Recipient:  m.Recipient,
			Total:      len(messages),
			Count:      i + 1, // Count starts at 1
			Reference:  ref,   // Reference should be the same for all messages
		}
	}
}

// Generate UDH header for concatenated message
func udh(ref int, total int, count int) string {
	return fmt.Sprintf("050003%02X%02X%02X", ref, total, count)
}

// Send message using the client
func (s *SmsService) send(message model.Message) {

	recipients := []string{strconv.Itoa(message.Recipient)}
	params := &messagebird.MessageParams{
		Type:        "binary",
		TypeDetails: messagebird.TypeDetails{"udh": udh(message.Reference, message.Total, message.Count)},
	}

	if s.client != nil {
		result, err := s.client.NewMessage(message.Originator, recipients, message.Message, params)

		if err != nil {
			log.Printf("Unable to send message, %v\n", err)

			for _, e := range result.Errors {
				log.Printf("Code: %d, Parameter: %s, Description: %s", e.Code, e.Parameter, e.Description)
			}
		}

	}
}

// Send messages at intervals
func startSend(timer <-chan time.Time, queue chan model.Message, done <-chan bool, s *SmsService) {
	log.Println("Sender started")
	for {
		select {
		case <-timer:

			select {
			case m := <-queue:
				log.Println("Sending " + m.String())

				s.send(m)

			case <-done:
				log.Println("Stopping sender")
				return
			}
		case <-done:
			log.Println("Stopping sender")
			return
		}
	}
}

// Split long message into chunks
func splitMessage(message string) (chunks []string) {
	length := len(message)
	count := (length / MaxLength) + 1

	for i := 0; i < count; i++ {
		to := int(math.Min(float64((i+1)*MaxLength), float64(length)))
		m := message[i*MaxLength : to]

		if m != "" {
			chunks = append(chunks, m)
		}

	}

	return chunks
}
