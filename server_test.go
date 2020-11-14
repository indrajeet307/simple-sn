package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)
func TestIndexPage(t *testing.T) {
	t.Run("Test the index page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		IndexHandler(response, request)
		rj := Response{}

		json.Unmarshal(response.Body.Bytes(), &rj)

		want := HuddleGreeting

		if rj.Message != want {
			t.Errorf("got `%q`, expected `%s`", rj.Message, want)
		}
	})
}

func TestUserOperations(t *testing.T) {
	t.Run("Test the index page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		IndexHandler(response, request)
		rj := Response{}

		json.Unmarshal(response.Body.Bytes(), &rj)

		want := HuddleGreeting

		if rj.Message != want {
			t.Errorf("got `%q`, expected `%s`", rj.Message, want)
		}
	})
	t.Run("Test new user add fails", func(t *testing.T) {
		jsonString := `{"Name":"test name", "meail":"a@b.com", "password":"testpassword"}`
		newUserBody, err := json.Marshal(jsonString)
		if err!=nil{
			t.Errorf("Failed to marshal %v", jsonString)
		}
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(newUserBody))
		request.Header.Add("content-type", "application/json")
		response := httptest.NewRecorder()

		AddNewUser(response, request)
		if response.Result().StatusCode == http.StatusOK {
			t.Errorf("Server was able to identify malformed request")
		}
	})
	t.Run("Test new user added", func(t *testing.T) {
		newUser := NewUserRequest{
			Name: "TestName",
			Email: "emailaddress@ex.com",
			Password: "testPassword",
		}
		response, err := addUser(newUser)
		if err != nil {
			t.Errorf("Failed to add new user %s", err.Error())
		}
		if response.Result().StatusCode != http.StatusOK {
			t.Errorf("Error occured on server")
		}
		rj := NewUserResponse{}

		json.Unmarshal(response.Body.Bytes(), &rj)

		if rj.Name != newUser.Name {
			t.Errorf("got `%q`, expected `%s`", rj.Name, newUser.Name)
		}
		if rj.Email != newUser.Email {
			t.Errorf("got `%q`, expected `%s`", rj.Email, newUser.Email)
		}
		if !rj.Active {
			t.Errorf("got `%t`, expected `%t`", rj.Active, true)
		}
	})
	t.Run("Test two users added", func(t *testing.T) {
		newUser := []NewUserRequest{
			{
				Name: "user1",
				Email: "user1@ex.com",
				Password: "user1pass",
			},
			{
				Name: "user2",
				Email: "user2@ex.com",
				Password: "user1pass",
			},
		}

		response1, err := addUser(newUser[0])
		if err != nil {
			t.Errorf("Failed to add new user %s", err.Error())
		}
		response2, err := addUser(newUser[1])
		if err != nil {
			t.Errorf("Failed to add new user %s", err.Error())
		}
		checkResponse := func (response *httptest.ResponseRecorder, id int){
			if response.Result().StatusCode != http.StatusOK {
				t.Errorf("Error occured on server")
			}
			rj := NewUserResponse{}

			json.Unmarshal(response.Body.Bytes(), &rj)

			if rj.ID != id {
				t.Errorf("got `%d`, expected `%d`", rj.ID, id)
			}
		}
		checkResponse(response1, 1)
		checkResponse(response2, 2)
	})
	t.Run("Test duplicate email id not allowed", func(t *testing.T) {
		newUser := []NewUserRequest{
			{
				Name: "user1",
				Email: "user1@ex.com",
				Password: "user1pass",
			},
			{
				Name: "user2",
				Email: "user1@ex.com",
				Password: "user2pass",
			},
		}

		_, err := addUser(newUser[0])
		if err != nil {
			t.Errorf("Failed to add new user %s", err.Error())
		}
		response, err := addUser(newUser[1])
		if response.Result().StatusCode == http.StatusOK {
			t.Errorf("Test should fail when adding two users with same email id")
		}
	})
}

func addUser(user NewUserRequest) (response *httptest.ResponseRecorder, err error){
	newUserBody, err := json.Marshal(user)
	if err!=nil{
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(newUserBody))
	request.Header.Add("content-type", "application/json")
	response = httptest.NewRecorder()

	AddNewUser(response, request)
	return response, nil
}
