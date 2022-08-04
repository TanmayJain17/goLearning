package controllers

import (
	"bytes"

	"encoding/json"
	"fmt"
	"go-fruit-cart/pkg/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAllProducts(t *testing.T) {

	r := gin.Default()
	r.GET("/products", GetAllProducts)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d%T\n", http.StatusOK, w.Code, w.Code)
	}
}

func TestAddProduct(t *testing.T) {
	cases := map[string]struct {
		status int

		body []byte

		authToken string

		resp gin.H
	}{

		"product added successfully": {

			status: http.StatusCreated,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken:/* adminToken, */ /* *adminTokenAd, */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			resp: gin.H{
				"message": "product added",
			},
		},
		"no access token": {

			status: http.StatusUnauthorized,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken: "",

			resp: gin.H{
				"error": "authorization header is missing",
			},
		},
		"must be admin to post": {

			status: http.StatusUnauthorized,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken:/* userToken, */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNTQzNH0.LKxvmQuW2LesB91Eckb0eqtWGarVUJ5Dsp-jLW83kEI",

			resp: gin.H{
				"error": "unauthorized",
			},
		},
		"access token invalid segments": {

			status: http.StatusUnauthorized,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken: "hello",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"invalid signature(access token)": {

			status: http.StatusUnauthorized,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken:/* adminToken + "somestring" */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWx",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"Name not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
					"name": "",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken:/* adminToken, */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Price not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
					"name": "product",
					"price":0,
					"description": "Great Qiality"
			}`),

			authToken:/* adminToken, */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Price not greater than 1": {

			status: http.StatusBadRequest,

			body: []byte(`{
					"name": "product",
					"price":1,
					"description": "Great Qiality"
			}`),

			authToken:/*  adminToken,  */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Description not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
					"name": "product",
					"price":50,
					"description": ""
			}`),

			authToken:/* adminToken, */ "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
	}
	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			token := v.authToken
			gin.SetMode(gin.TestMode)

			server := gin.New()
			server.Use(middleware.Auth())
			server.Handle(http.MethodPost, "/auth/products", AddProduct)

			httpServer := httptest.NewServer(server)

			requestURL := fmt.Sprintf("%s/auth/products", httpServer.URL)

			req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(v.body))

			if err != nil {

				t.Error("Unexpected Error", err.Error())

			}
			req.Header.Add("Authorization", fmt.Sprintf("%v", token))

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

func TestUpdateProduct(t *testing.T) {
	cases := map[string]struct {
		status int

		body []byte

		authToken string

		prodId string

		resp gin.H
	}{

		"product updated successfully": {

			status: http.StatusOK,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Quality"
			}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"message": "product added",
			},
		},
		"no access token": {

			status: http.StatusUnauthorized,

			body: []byte(`{
					"name": "Guvava",
					"price":70,
					"description": "Great Qiality"
			}`),

			authToken: "",
			prodId:    "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "authorization header is missing",
			},
		},

		"must be admin to update": {

			status: http.StatusUnauthorized,

			body: []byte(`{
						"name": "Guvava",
						"price":70,
						"description": "Great Qiality"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNTQzNH0.LKxvmQuW2LesB91Eckb0eqtWGarVUJ5Dsp-jLW83kEI",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "unauthorized",
			},
		},

		"access token invalid segments": {

			status: http.StatusUnauthorized,

			body: []byte(`{
						"name": "Guvava",
						"price":70,
						"description": "Great Qiality"
				}`),

			authToken: "hello",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"invalid signature(access token)": {

			status: http.StatusUnauthorized,

			body: []byte(`{
						"name": "Guvava",
						"price":70,
						"description": "Great Qiality"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWx",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"Name not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
						"name": "",
						"price":70,
						"description": "Great Qiality"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Price not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
						"name": "product",
						"price":0,
						"description": "Great Qiality"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Price not greater than 1": {

			status: http.StatusBadRequest,

			body: []byte(`{
						"name": "product",
						"price":1,
						"description": "Great Qiality"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Description not there": {

			status: http.StatusBadRequest,

			body: []byte(`{
						"name": "product",
						"price":50,
						"description": ""
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "023a0fff-fd9a-4bbc-9e89-9fb0764516f6",

			resp: gin.H{
				"error": "invalid input please check your input",
			},
		},
		"Wrong product id": {

			status: http.StatusNotFound,

			body: []byte(`{
						"name": "product",
						"price":50,
						"description": "description"
				}`),

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "197eac9d-ef2b-419c-abcd-76317891d0f7",

			resp: gin.H{
				"error": "product is not found",
			},
		},
	}
	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			token := v.authToken
			gin.SetMode(gin.TestMode)

			server := gin.New()
			server.Use(middleware.Auth())

			path := fmt.Sprintf("/auth/product/%v", v.prodId)

			server.Handle(http.MethodPut, path, UpdateProduct)

			/* httpServer := httptest.NewServer(server) */

			requestURL := fmt.Sprintf("http://127.0.0.1:8000/auth/product/%v", v.prodId)

			req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(v.body))
			/* ctx:=req.Context() */

			if err != nil {

				t.Error("Unexpected Error", err.Error())

			}

			req.Header.Add("Authorization", fmt.Sprintf("%v", token))

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

func TestDeleteProduct(t *testing.T) {
	cases := map[string]struct {
		status int

		authToken string

		prodId string

		resp gin.H
	}{

		"product deleted successfully": {

			status: http.StatusOK,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "9c6c14ba-ca19-4e7a-83d2-7b60a31593c3",

			resp: gin.H{
				"message": "successfully deleted",
			},
		},
		"no access token": {

			status: http.StatusUnauthorized,

			authToken: "",
			prodId:    "bf58c2c9-aec4-4865-89f2-f50fd15813b9",

			resp: gin.H{
				"error": "authorization header is missing",
			},
		},

		"must be admin to delete": {

			status: http.StatusUnauthorized,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNTQzNH0.LKxvmQuW2LesB91Eckb0eqtWGarVUJ5Dsp-jLW83kEI",

			prodId: "197eac9d-ef2b-419c-af9d-76317891d0f7",

			resp: gin.H{
				"error": "unauthorized",
			},
		},

		"access token invalid segments": {

			status: http.StatusUnauthorized,

			authToken: "hello",

			prodId: "197eac9d-ef2b-419c-af9d-76317891d0f7",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"invalid signature(access token)": {

			status: http.StatusUnauthorized,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWx",

			prodId: "197eac9d-ef2b-419c-af9d-76317891d0f7",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"Wrong product id": {

			status: http.StatusNotFound,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTUzNTh9.96P1pwJLCYEWVNIZMrdCiGqU-HboyGlfNI6GelWxjLo",

			prodId: "197eac9d-ef2b-419c-abcd-76317891d0f7",

			resp: gin.H{
				"error": "product is not found",
			},
		},
	}
	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			token := v.authToken
			gin.SetMode(gin.TestMode)

			server := gin.New()
			server.Use(middleware.Auth())

			path := fmt.Sprintf("/auth/product/%v", v.prodId)
			/* fmt.Println(path) */

			server.Handle(http.MethodDelete, path, DeleteProduct)

			/* httpServer := httptest.NewServer(server) */

			requestURL := fmt.Sprintf("http://127.0.0.1:8000/auth/product/%v", v.prodId)
			/* fmt.Println(requestURL)
			fmt.Println(httpServer.URL) */
			req, err := http.NewRequest(http.MethodDelete, requestURL, nil /* bytes.NewBuffer(v.body) */)
			/* ctx:=req.Context() */

			if err != nil {

				t.Error("Unexpected Error", err.Error())

			}

			req.Header.Add("Authorization", fmt.Sprintf("%v", token))

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
