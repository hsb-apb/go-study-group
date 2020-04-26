package model

type Request struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
}

type Response struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
