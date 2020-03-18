package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type VmsListRequest struct {
	offset int
}

func MakeVmsListRequestPayload(offset int) VmsListRequest {
	return VmsListRequest{offset: offset}
}

func GetVmsLists(url string, username string, password string, payload VmsListRequest) *http.Response {
	jsonData, _ := json.Marshal(payload)
	buff := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
