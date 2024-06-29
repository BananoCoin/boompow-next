package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/bananocoin/boompow/libs/models"
	"k8s.io/klog/v2"
)

type RPCClient struct {
	Url        string
	httpClient *http.Client
}

func NewRPCClient(url string) *RPCClient {
	return &RPCClient{
		Url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 30, // Set a timeout for all requests
		},
	}
}

type SendResponse struct {
	Block string `json:"block"`
}

// Base request
func (client RPCClient) makeRequest(request interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(request)
	// HTTP post
	resp, err := client.httpClient.Post(client.Url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		klog.Errorf("Error making RPC request %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Try to decode+deserialize
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Error decoding response body %s", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		klog.Errorf("Received non-200 response: %s", body)
		return nil, errors.New("non-200 response received")
	}

	return body, nil
}

// send
func (client RPCClient) MakeSendRequest(request models.SendRequest) (*SendResponse, error) {
	response, err := client.makeRequest(request)
	if err != nil {
		klog.Errorf("Error making request %s", err)
		return nil, err
	}
	// Try to decode+deserialize
	var sendResponse SendResponse
	err = json.Unmarshal(response, &sendResponse)
	if err != nil {
		klog.Errorf("Error unmarshaling response %s, %s", string(response), err)
		return nil, errors.New("Error")
	}
	return &sendResponse, nil
}
