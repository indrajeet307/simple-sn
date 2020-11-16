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
}

type CommentReplyRequest struct {
	ID       int64  `json:"id"`
	FromUser int64  `json:"from_user"`
	Body     string `json:"body"`
}

type NewCommentResponse struct {
	ID int64 `json:"id"`
}

type WallCommentsResponse struct {
	Comments []NewCommentRequest
}

type CommentReactionRequest struct {
	ReactionID int64 `json:"reaction_id"`
}

type CommentReactionResponse struct {
	CommentID  int64 `json:"comment_id"`
	ReactionID int64 `json:"reaction_id"`
	Count      int64 `json:"count"`
}

type CommentListReactions struct {
	Reactions []CommentReactionResponse `json:"reactions"`
}

type ReactionRequest struct {
	ID   int64  `json:"reaction_id"`
	Name string `json:"name"`
}
type ReactionResponse struct {
	Reactions []ReactionRequest `json:"reactions"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
