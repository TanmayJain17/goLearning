package controllers

import (
	"encoding/json"
	"fmt"
	"go-fruit-cart/pkg/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetCart(t *testing.T) {
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06MPM"
	r := gin.Default()
	r.Use(middleware.Auth())
	r.GET("/auth/cart", GetCart)

	req, err := http.NewRequest(http.MethodGet, "/auth/cart", nil)
	req.Header.Add("Authorization", fmt.Sprintf("%v", token))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	/* fmt.Println(w.Body)

	fmt.Printf("%v %T", w.Code, w.Code)
	fmt.Printf("%v %T", http.StatusOK, http.StatusOK) */

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d%T\n", http.StatusOK, w.Code, w.Code)
	}
}

func TestAddToCart(t *testing.T) {
	cases := map[string]struct {
		status int

		authToken string

		prodId string

		resp gin.H
	}{

		"product added successfully": {

			status: http.StatusOK,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06MPM",

			prodId: "aa7faf53-6001-4b96-a02a-67708d3a72bb",

			resp: gin.H{
				"message": "product added",
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

		"must be user to update": {

			status: http.StatusUnauthorized,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTY5MTN9.l3wiHHjy0gDFFD9s1jaLhgXBFG_s_kHvpiz0Qu8iqcE",

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

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06",

			prodId: "197eac9d-ef2b-419c-af9d-76317891d0f7",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"Wrong product id": {

			status: http.StatusBadRequest,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06MPM",

			prodId: "69643226-78dd-4cd1-abca-b3a95345be",

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

			path := fmt.Sprintf("/auth/products/%v/cart", v.prodId)
			/* fmt.Println(path) */

			server.Handle(http.MethodPut, path, AddToCart)

			/* httpServer := httptest.NewServer(server) */

			requestURL := fmt.Sprintf("http://127.0.0.1:8000/auth/products/%v/cart", v.prodId)
			/* fmt.Println(requestURL)
			fmt.Println(httpServer.URL) */
			req, err := http.NewRequest(http.MethodPut, requestURL, nil /* bytes.NewBuffer(v.body) */)
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

func TestRemoveFromCart(t *testing.T) {
	cases := map[string]struct {
		status int

		authToken string

		prodId string

		resp gin.H
	}{

		"product removed successfully": {

			status: http.StatusOK,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06MPM",

			prodId: "7b0f0075-ec62-4171-9265-4cea81b38c64",

			resp: gin.H{
				"message": "product removed",
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

		"must be user to update": {

			status: http.StatusUnauthorized,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSYWh1bCIsImxhc3RuYW1lIjoiRHJhdmlkIiwiZW1haWwiOiJyZEBmYy5jb20iLCJleHAiOjE2NTk0MTY5MTN9.l3wiHHjy0gDFFD9s1jaLhgXBFG_s_kHvpiz0Qu8iqcE",

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

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06",

			prodId: "197eac9d-ef2b-419c-af9d-76317891d0f7",

			resp: gin.H{
				"error": "auth token is invalid",
			},
		},
		"Wrong product id": {

			status: http.StatusBadRequest,

			authToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdG5hbWUiOiJSb2hpdCIsImxhc3RuYW1lIjoiTWVocmEiLCJlbWFpbCI6InJtQGZjLmNvbSIsImV4cCI6MTY1OTQxNjkwNH0.FhpqlVn6mR3I8x915iKfZfKS63CfuY_2JkqzKX06MPM",

			prodId: "69643226-78dd-4cd1-abca-b3a95345be6b",

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

			path := fmt.Sprintf("/auth/products/%v/cart", v.prodId)
			/* fmt.Println(path) */

			server.Handle(http.MethodDelete, path, RemoveFromCart)

			/* httpServer := httptest.NewServer(server) */

			requestURL := fmt.Sprintf("http://127.0.0.1:8000/auth/products/%v/cart", v.prodId)
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
