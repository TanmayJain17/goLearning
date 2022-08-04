package controllers

import (
	"bytes"
	/* "context" */
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	/* "github.com/google/uuid" */)

func TestNewUser(t *testing.T) {
	cases := map[string]struct {
		status int

		body []byte

		resp gin.H

		dbErr error
	}{

		"register successfull": {

			status: http.StatusCreated,

			body: []byte(`{
	
					"email": "test-u225@test.com",
	
					"password": "test-password",
	
					"firstName": "test",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{
				"message": "SignUp Successfull",
			},
		},

		"email format is invalid": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "test",
	
					"password": "test-password",
	
					"firstName": "test",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},
		"email field is empty": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "",
	
					"password": "test-password",
	
					"firstName": "test",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},

		"password is too short": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "test-user@test.com",
	
					"password": "te",
	
					"firstName": "test",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},

		"first name field is invalid": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "test-user@test.com",
	
					"password": "test-password",
	
					"firstName": "  ",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},
		"fisrt name field is empty": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "abc@g.com",
	
					"password": "test-password",
	
					"firstName": "",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},
		"last name field is empty": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "abc@g.com",
	
					"password": "test-password",
	
					"firstName": "ewferf",
	
					"lastName": ""
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},

		"last name is invalid": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "test-user@test.com",
	
					"password": "test-password",
	
					"firstName": "test",
	
					"lastName": "  "
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},

		"email already exists": {

			status: http.StatusUnauthorized,

			body: []byte(`{
	
					"email":"bj@fc.com",
	
					"password": "test-password",
	
					"firstName": "test",
	
					"lastName": "user"
	
				}`),

			resp: gin.H{

				"error": "user already exists",
			},
		},
	}
	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			gin.SetMode(gin.TestMode)

			server := gin.New()

			server.Handle(http.MethodPost, "/signup", NewUser)

			httpServer := httptest.NewServer(server)

			requestURL := fmt.Sprintf("%s/signup", httpServer.URL)

			req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(v.body))

			if err != nil {

				t.Error("Unexpected Error", err.Error())

			}

			client := http.Client{}

			res, err := client.Do(req)

			if err != nil {

				t.Error("unexpected error: ", err)

			}

			if status := res.StatusCode; status != v.status {

				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)

			}

			if res.StatusCode != http.StatusOK {

				body, err := ioutil.ReadAll(res.Body)

				if err != nil {

					t.Error("unexpected error: ", err)

				}

				var got gin.H

				err = json.Unmarshal(body, &got)

				if err != nil {

					t.Fatal(err)

				}

				if fmt.Sprint(v.resp) != fmt.Sprint(got) {

					t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)

				}

			}

		})
	}
}

func TestLogin(t *testing.T) {
	cases := map[string]struct {
		status int

		body []byte

		resp gin.H

		dbErr error
	}{

		"login successfull": {

			status: http.StatusOK,

			body: []byte(`{
	
					"email":"bj@fc.com",
					"password":"bj@1234"
	
					
	
				}`),

			resp: gin.H{
				"message": "successfully logged in",
			},
		},

		"email format is invalid": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "test",
	
					"password": "test-password"
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},
		"email field is empty": {

			status: http.StatusBadRequest,

			body: []byte(`{
	
					"email": "",
	
					"password": "test-password"
	
					
	
				}`),

			resp: gin.H{

				"error": "invalid input please check your input",
			},
		},
		"email does not exists": {

			status: http.StatusNotFound,

			body: []byte(`{
	
					"email":"bj@mf.com",
	
					"password": "test-password"
	
					
	
				}`),

			resp: gin.H{

				"error": "user does not exist",
			},
		},
		"wrong password": {

			status: http.StatusUnauthorized,

			body: []byte(`{
	
					"email":"bj@fc.com",
	
					"password": "bj1234"
	
					
	
				}`),

			resp: gin.H{

				"error": "invalid credentials",
			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			gin.SetMode(gin.TestMode)

			server := gin.New()

			server.Handle(http.MethodPost, "/login", UserLogin)

			httpServer := httptest.NewServer(server)

			requestURL := fmt.Sprintf("%s/login", httpServer.URL)

			req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(v.body))

			if err != nil {

				t.Error("Unexpected Error", err.Error())

			}

			client := http.Client{}

			res, err := client.Do(req)

			if err != nil {

				t.Error("unexpected error: ", err)

			}

			if status := res.StatusCode; status != v.status {

				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.status)

			}

			if res.StatusCode != http.StatusOK {

				body, err := ioutil.ReadAll(res.Body)

				if err != nil {

					t.Error("unexpected error: ", err)

				}

				var got gin.H

				err = json.Unmarshal(body, &got)

				if err != nil {

					t.Fatal(err)

				}

				if fmt.Sprint(v.resp) != fmt.Sprint(got) {

					t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)

				}

			}

		})
	}
}
