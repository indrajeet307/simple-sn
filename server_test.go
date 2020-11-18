package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"fmt"

	"github.com/gorilla/mux"
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
		defer NewDB()

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
		defer NewDB()

		jsonString := `{"Name":"test name", "meail":"a@b.com", "password":"testpassword"}`
		newUserBody, err := json.Marshal(jsonString)
		if err != nil {
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
		defer NewDB()

		newUser := NewUserRequest{
			Name:     "TestName",
			Email:    "emailaddress@ex.com",
			Password: "testPassword",
		}
		response, err := addUser(newUser)
		if err != nil {
			t.Fatalf("Failed to add new user %s", err.Error())
		}
		if response.Result().StatusCode != http.StatusOK {
			t.Logf("%s", string(response.Body.Bytes()))
			t.Fatalf("Api return code no ok")
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
		defer NewDB()

		users := []NewUserRequest{
			{
				Name:     "user1",
				Email:    "user1@ex.com",
				Password: "user1pass",
			},
			{
				Name:     "user2",
				Email:    "user2@ex.com",
				Password: "user1pass",
			},
		}

		addUserAndCheckResponse := func(nu NewUserRequest, id int64) {
			response, err := addUser(nu)
			if err != nil {
				t.Fatalf("Failed to add new user %s", err.Error())
			}
			if response.Result().StatusCode != http.StatusOK {
				t.Fatalf("Api return code no ok")
			}
			rj := NewUserResponse{}

			json.Unmarshal(response.Body.Bytes(), &rj)

			if rj.ID != id {
				t.Errorf("got `%d`, expected `%d`", rj.ID, id)
			}
		}
		addUserAndCheckResponse(users[0], 1)
		addUserAndCheckResponse(users[1], 2)
	})
	t.Run("Test duplicate email id not allowed", func(t *testing.T) {
		defer NewDB()

		newUser := []NewUserRequest{
			{
				Name:     "user1",
				Email:    "user1@ex.com",
				Password: "user1pass",
			},
			{
				Name:     "user2",
				Email:    "user1@ex.com",
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

func addUser(user NewUserRequest) (response *httptest.ResponseRecorder, err error) {
	newUserBody, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(newUserBody))
	request.Header.Add("content-type", "application/json")
	response = httptest.NewRecorder()

	AddNewUser(response, request)
	return response, nil
}

func TestWallOperations(t *testing.T) {
	t.Run("Test can add to own wall", func(t *testing.T) {
		db = GetDB()
		defer NewDB()

		db.AddUser( &NewUserRequest{
			Name: "test",
			Email: "e@b.com",
			Password:"password",
		})

		newComment := NewCommentRequest{
			ToUser:   1,
			FromUser: 1,
			Body:     "Some mildly intresting body",
		}

		response, err := addComment(newComment, "1")
		if err != nil {
			t.Fatalf("Unable to add comment")
		}

		rj := NewCommentResponse{}

		json.Unmarshal(response.Body.Bytes(), &rj)

		want := int64(1)

		if rj.ID != want {
			t.Errorf("got `%d`, expected `%d`", rj.ID, want)
		}
	})
	t.Run("Test can add multiple comments to wall", func(t *testing.T) {
		db = GetDB()
		defer NewDB()

		db.AddUser( &NewUserRequest{
			Name: "user1",
			Email: "user1@b.com",
			Password:"password",
		})
		db.AddUser( &NewUserRequest{
			Name: "user2",
			Email: "user2@b.com",
			Password:"password",
		})

		newComments := []NewCommentRequest{
			{
				FromUser: 1,
				Body:     "Some intresting body",
			},
			{
				FromUser: 2,
				Body:     "Some intresting body 2",
			},
		}

		response, err := addComment(newComments[0], "1")
		if err != nil {
			t.Fatalf("Unable to add comment")
			return
		}
		rj := NewCommentResponse{}
		json.Unmarshal(response.Body.Bytes(), &rj)
		if rj.ID != 1 {
			t.Errorf("got `%d`, expected `%d`", rj.ID, 1)
		}

		response, err = addComment(newComments[1], "1")
		if err != nil {
			t.Fatalf("Unable to add comment")
			return
		}
		rj = NewCommentResponse{}
		json.Unmarshal(response.Body.Bytes(), &rj)
		if rj.ID != 2 {
			t.Errorf("got `%d`, expected `%d`", rj.ID, 2)
		}

	})
	t.Run("Test can fetch wall", func(t *testing.T) {
		db = GetDB()
		defer NewDB()

		db.AddUser( &NewUserRequest{
			Name: "user1",
			Email: "user1@b.com",
			Password:"password",
		})
		db.AddUser( &NewUserRequest{
			Name: "user2",
			Email: "user2@b.com",
			Password:"password",
		})

		newComments := []NewCommentRequest{
			{
				FromUser: 1,
				Body:     "Some intresting body",
			},
			{
				FromUser: 2,
				Body:     "Some intresting body 2",
			},
		}

		addComment(newComments[0], "1")
		addComment(newComments[1], "1")

		request, _ := http.NewRequest(http.MethodGet, "/wall/", nil)
		request = mux.SetURLVars(request, map[string]string{
			"userID": "1",
		})
		response := httptest.NewRecorder()

		GetUserWall(response, request)

		wcr := WallCommentsResponse{}
		json.Unmarshal(response.Body.Bytes(), &wcr)

		if len(wcr.Comments) != 2 {
			t.Errorf("Failed To Fetch Wall Comments")
		}

		for _, comment := range wcr.Comments {
			if comment.ToUser == 3 {
				t.Errorf("Invalid Response, Should Fetch Entries From User 1 Only")
			}
		}

	})
}

func addComment(comment NewCommentRequest, uid string) (response *httptest.ResponseRecorder, err error) {
	requestBody, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, "/wall", bytes.NewBuffer(requestBody))
	request = mux.SetURLVars(request, map[string]string{
		"userID": uid,
	})
	response = httptest.NewRecorder()
	AddToWall(response, request)
	return response, nil
}

func addCommentReply(comment CommentReplyRequest, cid int64) (response *httptest.ResponseRecorder, err error) {
	requestBody, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(requestBody))
	request = mux.SetURLVars(request, map[string]string{
		"commentID": strconv.FormatInt(cid, 10),
	})
	response = httptest.NewRecorder()
	AddCommentReply(response, request)
	return response, nil

}

func TestCommentOperation(t *testing.T) {
	t.Run("test can add comment", func(t *testing.T) {
		db = GetDB()
		defer NewDB()

		db.AddUser( &NewUserRequest{
			Name: "user1",
			Email: "user1@b.com",
			Password:"password",
		})
		db.AddUser( &NewUserRequest{
			Name: "user2",
			Email: "user2@b.com",
			Password:"password",
		})

		newComment := NewCommentRequest{
			FromUser: 1,
			Body:     "Some intresting body",
		}
		response, err := addComment(newComment, "1")
		if err != nil {
			t.Errorf("Failed to add a new comment")
		}

		commentResponse := NewCommentResponse{}
		json.Unmarshal(response.Body.Bytes(), &commentResponse)

		commentReply := CommentReplyRequest{
			FromUser: 2,
			Body:     "This is really a nice comment",
		}

		t.Log(commentResponse.ID)

		response, err = addCommentReply(commentReply, commentResponse.ID)
		if err != nil {
			t.Errorf("Failed to add a new comment")
			return
		}

		commentReplyResponse := NewCommentResponse{}
		json.Unmarshal(response.Body.Bytes(), &commentReplyResponse)

		if commentReplyResponse.ID != 2 {
			t.Log(commentReplyResponse.ID)
			t.Errorf("Unable to add comment reply to existing comment")
		}

	})
}

func addCommentReaction(cid int64, rr CommentReactionRequest) (response *CommentReactionResponse, err error) {
	requestBody, err := json.Marshal(rr)
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/reactions/%d", cid) , bytes.NewBuffer(requestBody))
	recorder := httptest.NewRecorder()
	AddCommentReaction(recorder, request)

	response = &CommentReactionResponse{}
	json.Unmarshal(recorder.Body.Bytes(), response)
	return
}

func getCommentReactions(cid int) (response *httptest.ResponseRecorder, err error) {
	request, _ := http.NewRequest(http.MethodGet, "/reactions", nil)
	request = mux.SetURLVars(request, map[string]string{
		"commentID": strconv.Itoa(cid),
	})
	response = httptest.NewRecorder()
	GetCommentReaction(response, request)
	return response, nil
}

func TestReactionOperation(t *testing.T) {
	t.Run("test can add reaction to comment", func(t *testing.T) {
		db = GetDB()
		defer NewDB()

		db.AddUser( &NewUserRequest{
			Name: "test",
			Email: "e@b.com",
			Password:"password",
		})
		db.AddReaction( &ReactionRequest{
			Name: "testreact",
		})

		newComment := NewCommentRequest{
			FromUser: 1,
			Body:     "Some intresting body",
		}
		response, err := addComment(newComment, "1")
		if err != nil {
			t.Errorf("Failed to add a new comment")
		}

		commentResponse := NewCommentResponse{}
		json.Unmarshal(response.Body.Bytes(), &commentResponse)

		commentReaction := CommentReactionRequest{
			ReactionID: 1,
		}

		resp, err := addCommentReaction(commentResponse.ID, commentReaction)
		if err != nil {
			t.Fatalf("Failed to add a new reaction")
		}

		if resp.CommentID != commentResponse.ID {
			t.Logf("%d  %d", resp.CommentID, commentResponse.ID)
			t.Fatal("Unable to add comment reply to existing comment")
		}

		response, err = getCommentReactions(int(commentResponse.ID))
		if err != nil {
			t.Fatalf("Failed to fetch reactions for comment")
			return
		}

		listReactions := CommentListReactions{}
		json.Unmarshal(response.Body.Bytes(), &listReactions)

		if len(listReactions.Reactions) != 1 {
			t.Fatalf("Unable to read added reaction properly")
		}
		if listReactions.Reactions[0].Count != 1 {
			t.Errorf("Reactions count not set properly")
		}

		resp, err = addCommentReaction(commentResponse.ID, commentReaction)
		if err != nil {
			t.Fatalf("Failed to add a new reaction")
		}
		resp, err = addCommentReaction(commentResponse.ID, commentReaction)
		if err != nil {
			t.Fatalf("Failed to add a new reaction")
		}

		response, err = getCommentReactions(int(commentResponse.ID))
		if err != nil {
			t.Fatalf("Failed to fetch reactions for comment")
			return
		}

		listReactions = CommentListReactions{}
		json.Unmarshal(response.Body.Bytes(), &listReactions)

		if len(listReactions.Reactions) != 1 {
			t.Errorf("Unable to read added reaction properly")
		}
		if listReactions.Reactions[0].Count != 3 {
			t.Errorf("Reactions count not set properly")
		}

		commentReaction.ReactionID = 2
		resp, err = addCommentReaction(commentResponse.ID, commentReaction)
		if err != nil {
			t.Fatalf("Failed to add a new reaction")
		}
		response, err = getCommentReactions(int(commentResponse.ID))
		if err != nil {
			t.Fatalf("Failed to fetch reactions for comment")
			return
		}

		listReactions = CommentListReactions{}
		json.Unmarshal(response.Body.Bytes(), &listReactions)

		if len(listReactions.Reactions) != 2 {
			t.Errorf("Unable to read added reaction properly")
		}
		if listReactions.Reactions[1].Count != 1 {
			t.Errorf("Reactions count not set properly")
		}

	})
}
