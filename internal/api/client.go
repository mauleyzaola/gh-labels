package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mauleyzaola/gh-labels/internal/types"
)

type Client struct {
	client *http.Client
	token  string
}

func New(token string) (*Client, error) {
	return &Client{
		client: &http.Client{},
		token:  token,
	}, nil
}

func (x *Client) genResponse(method, endpoint string, payload []byte, statusCode int) (*http.Response, error) {
	const baseUrl = "https://api.github.com"
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	uri := baseUrl + "/" + endpoint
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		log.Println("[ERROR]", err)
		return nil, err
	}
	req.Header.Set("Authorization", "token "+x.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := x.client.Do(req)
	if err != nil {
		return nil, err
	}
	if expected, actual := statusCode, res.StatusCode; expected != actual {
		return nil, fmt.Errorf("unexpected response status code: expected %d, actual %d", expected, actual)
	}
	return res, nil
}

func (x *Client) List(info types.RepoInfo) ([]types.Label, error) {
	res, err := x.genResponse(http.MethodGet,
		"repos/"+info.Username+"/"+info.Repository+"/labels",
		nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var result []types.Label
	defer func() { _ = res.Body.Close() }()
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (x *Client) Create(dst types.RepoInfo, label types.Label) error {
	payload, err := json.Marshal(label)
	if err != nil {
		return err
	}
	res, err := x.genResponse(http.MethodPost,
		"repos/"+dst.Username+"/"+dst.Repository+"/labels",
		payload, http.StatusCreated)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	// TODO: report this bug to GH, description is not getting there
	return json.NewDecoder(res.Body).Decode(&label)
}

func (x *Client) Delete(dst types.RepoInfo, labelName string) error {
	res, err := x.genResponse(http.MethodDelete,
		"repos/"+dst.Username+"/"+dst.Repository+"/labels/"+labelName,
		nil, http.StatusNoContent)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	return nil
}
