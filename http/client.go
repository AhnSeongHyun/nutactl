package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type VmsListRequest struct {
	Offset int `json:"offset"`
}

func MakeVmsListRequestPayload(offset int) VmsListRequest {
	return VmsListRequest{Offset: offset}
}

func GetVmsLists(url string, username string, password string, payload VmsListRequest) *http.Response {
	pbytes, _ := json.Marshal(payload)
	buff := bytes.NewBuffer(pbytes)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
