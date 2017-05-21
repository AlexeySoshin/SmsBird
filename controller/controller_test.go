package controller

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendMessage(t *testing.T) {
	c := NewController("")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.SendSms(w, r)
	})

	Convey("Don't allow GET", t, func() {
		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusMethodNotAllowed)
	})

	Convey("Don't allow POST empty body", t, func() {
		req := httptest.NewRequest("POST", "/", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)

	})

	Convey("Allow POST with body", t, func() {
		body := `{"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}`
		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)

	})

	Convey("Should not allow", t, func() {

		Convey("Empty message", func() {
			body := `{"recipient":31612345678,"originator":"MessageBird","message":"               "}`
			req := httptest.NewRequest("POST", "/", strings.NewReader(body))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("Empty recipient", func() {
			body := `{"recipient":"","originator":"MessageBird","message":"This is a test message."}`
			req := httptest.NewRequest("POST", "/", strings.NewReader(body))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("Bad recepient", func() {
			body := `{"recipient":123,"originator":"MessageBird","message":"This is a test message."}`
			req := httptest.NewRequest("POST", "/", strings.NewReader(body))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

	})
}
