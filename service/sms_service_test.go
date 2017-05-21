package service

import (
	"fmt"
	"github.com/alexeysoshin/SmsBird/model"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSplitMessageShort(t *testing.T) {

	message := "abc"

	result := splitMessage(message)

	Convey("Should return same message", t, func() {

		So(len(result), ShouldEqual, 1)
		So(result[0], ShouldEqual, message)
	})
}

func TestSplitMessageExact(t *testing.T) {

	message := ""
	for i := 0; i < MaxLength; i++ {
		message += "a"
	}

	result := splitMessage(message)

	Convey("Should return same message", t, func() {

		So(len(result), ShouldEqual, 1)
		So(result[0], ShouldEqual, message)
	})
}

func TestSplitMessageLong(t *testing.T) {

	message := ""
	for i := 0; i < MaxLength+2; i++ {
		if i%2 == 0 {
			message += "a"
		} else {
			message += "b"
		}

	}

	result := splitMessage(message)

	Convey("Should split message if message is too long", t, func() {

		So(len(result), ShouldEqual, 2)
		So(len(result[0]), ShouldEqual, 160)
		So(result[1], ShouldEqual, "ab")
	})
}

func TestUdh(t *testing.T) {
	Convey("UDH for concatenated SMS", t, func() {

		u := udh(247, 3, 2)

		So(u, ShouldEqual, "050003F70302")
	})
}

func TestSendMessage(t *testing.T) {

	Convey("Should send messages in queue", t, func() {

		done := make(chan bool)
		Convey("Should empty queue if possible", func() {

			queue := make(chan model.Message, 100)
			messageCount := 5

			for i := 0; i < messageCount; i++ {
				queue <- model.Message{Message: fmt.Sprintf("message%d", i)}
			}
			So(len(queue), ShouldEqual, messageCount)
			timer := time.NewTicker(time.Millisecond * 5).C

			go startSend(timer, queue, done, nil)

			timeoutTimer := time.NewTimer(time.Millisecond * 50).C

			select {
			case <-timeoutTimer:
				done <- true
			}

			So(len(queue), ShouldEqual, 0)
		})

		Convey("Should stop even if queue is not empty", func() {
			queue := make(chan model.Message, 100)
			messageCount := 5

			for i := 0; i < messageCount; i++ {
				queue <- model.Message{Message: fmt.Sprintf("message%d", i)}
			}
			So(len(queue), ShouldEqual, messageCount)
			timer := time.NewTicker(time.Millisecond * 5).C

			go startSend(timer, queue, done, nil)

			timeoutTimer := time.NewTimer(time.Millisecond * 8).C

			select {
			case <-timeoutTimer:
				done <- true
			}

			So(len(queue), ShouldBeGreaterThan, 0)
		})
	})
}
