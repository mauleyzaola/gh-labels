package main

import "net/http"

const (
	GHBaseURL = "https://api.github.com"
)

type Label struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type APIClient struct {
	baseURL string
	client  *http.Client
	token   string
}
