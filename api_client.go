package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// based on GH documentation: https://developer.github.com/v3/issues/labels/

func NewAPIClient(timeout *time.Duration) (*APIClient, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("cannot resolve environment variable: TOKEN")
	}
	return &APIClient{
		baseURL: GHBaseURL,
		client:  NewHTTPClient(timeout),
		token:   token,
	}, nil
}

func NewHTTPClient(timeout *time.Duration) *http.Client {
	if timeout != nil {
		return &http.Client{
			Timeout: *timeout,
		}
	} else {
		return &http.Client{
			Timeout: time.Second * 30,
		}
	}
}

func (ac *APIClient) setAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "token "+ac.token)
}

func (ac *APIClient) createResponse(req *http.Request) (*http.Response, error) {
	ac.setAuthHeader(req)
	res, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ac *APIClient) debugBody(res *http.Response) {
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println(string(data))
}

// List returns all the labels from a given repository
func (ac *APIClient) LabelList(username, repository string) ([]Label, error) {
	address := fmt.Sprintf("%s/repos/%s/%s/labels", ac.baseURL, username, repository)
	uri, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	var result []Label
	res, err := ac.createResponse(req)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (ac *APIClient) LabelPost(username, repository string, label *Label) error {
	if label == nil {
		return errors.New("nil parameter: label")
	}

	address := fmt.Sprintf("%s/repos/%s/%s/labels", ac.baseURL, username, repository)
	data, err := json.Marshal(label)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, address, body)
	if err != nil {
		return err
	}
	res, err := ac.createResponse(req)
	if err != nil {
		return err
	}
	if expected, actual := http.StatusCreated, res.StatusCode; expected != actual {
		ac.debugBody(res)
		return fmt.Errorf("expected status code:%d got: %d", expected, actual)
	}

	// TODO: report this bug to GH, description is not getting there
	return json.NewDecoder(res.Body).Decode(label)
}

func (ac *APIClient) LabelDelete(username, repository, labelName string) error {
	address := fmt.Sprintf("%s/repos/%s/%s/labels/%s", ac.baseURL, username, repository, labelName)
	uri, err := url.Parse(address)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	if err != nil {
		return err
	}
	res, err := ac.createResponse(req)
	if err != nil {
		return err
	}
	if expected, actual := http.StatusNoContent, res.StatusCode; expected != actual {
		ac.debugBody(res)
		return fmt.Errorf("expected status code:%d got: %d", expected, actual)
	}
	return nil
}
