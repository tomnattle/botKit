package dfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"io/ioutil"
	"net/http"
)

var cfg *config.DfsConfig

func init() {
	cfg = config.GetConfig().Dfs
	if cfg == nil {
		panic("Dfs config is nil")
	}
}

type Response struct {
	ErrNo   int64       `json:"err_no"`
	ErrMsg  string      `json:"err_msg"`
	Results interface{} `json:"results"`
}

type Result struct {
	Header   commonHTTP.Header `json:"header"`
	Response Response          `json:"response"`
}

type ReadRequest struct {
	GroupName   string `json:"groupname"`
	FileName    string `json:"filename"`
	OffSet      int    `json:"offset"`
	Length      int    `json:"length"`
	ContentType string `json:"contentType"`
}

func Read(r ReadRequest) (*Response, error) {
	if r.GroupName == "" || r.FileName == "" {
		return nil, fmt.Errorf("params error")
	}

	req := &ReadRequest{}
	reqBody := commonHTTP.MakeReq(&req)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "download"
	reqBody.R.P = r

	post, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			continue
		}
		if response.StatusCode != 200 {
			continue
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		result := Result{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}

		if result.Response.ErrMsg != "" {
			return nil, fmt.Errorf("read error: %s", result.Response.ErrMsg)
		}

		return &result.Response, nil
	}

	return nil, fmt.Errorf("read error")
}

type UploadRequest struct {
	Content     string `json:"content"`
	Ext         string `json:"ext"`
	ContentType string `json:"contentType"`
}

func Write(w UploadRequest) (*Response, error) {
	if w.Content == "" {
		return nil, fmt.Errorf("params error")
	}

	if w.ContentType == "" {
		w.ContentType = "json"
	}
	if w.Ext == "" {
		w.Ext = "txt"
	}

	req := &ReadRequest{}
	reqBody := commonHTTP.MakeReq(&req)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "upload"
	reqBody.R.P = w

	post, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			continue
		}
		if response.StatusCode != 200 {
			continue
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		result := Result{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}

		if result.Response.ErrMsg != "" {
			return nil, fmt.Errorf("read error: %s", result.Response.ErrMsg)
		}

		return &result.Response, nil
	}
	return nil, fmt.Errorf("upload error")
}

func Del(d ReadRequest) (bool, error) {
	if d.GroupName == "" || d.FileName == "" {
		return false, fmt.Errorf("params error")
	}
	req := &ReadRequest{}
	reqBody := commonHTTP.MakeReq(&req)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "del"
	reqBody.R.P = d

	post, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			continue
		}
		if response.StatusCode != 200 {
			continue
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		result := Result{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			return false, err
		}

		if result.Response.ErrMsg != "" {
			return false, fmt.Errorf("read error: %s", result.Response.ErrMsg)
		}

		return true, nil
	}

	return false, fmt.Errorf("del error")
}
