package main

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type NewUserRequest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   string `json:"-"`
}

type NewUserResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type NewCommentRequest struct {
	ID       int64  `json:"id"`
	ToUser   int64  `json:"to_user"`
	FromUser int64  `json:"from_user"`
	Body     string `json:"body"`
	ParentID int64  `json:"parent_id"`
}

type NewCommentResponse struct {
	ID int64 `json:"id"`
}

type WallCommentsResponse struct {
	Comments []NewCommentRequest
}

type ReactionRequest struct {
	CommentID  int64 `json:"comment_id"`
	ReactionID int64 `json:"reaction_id"`
}

type ReactionResponse struct {
	CommentID int64 `json:"comment_id"`
}

type ListReactions struct {
	Reactions []ReactionRequest `json:"reactions"`
}
