package dfs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hsyan2008/hfw2/curl"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
)

var cfg *config.DfsConfig
var req commonHTTP.Request

func init() {
	cfg = config.GetConfig().Dfs
	if cfg == nil {
		panic("Dfs config is nil")
	}
	req.H = commonHTTP.Header{
		AppID:    config.GetConfig().AppID,
		LogID:    config.GetConfig().Environment,
		Provider: config.GetConfig().AppName,
	}

	req.R.C = "Dfs"
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

	req.R.M = "download"
	req.R.P = r

	post, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json marshal error")
	}

	result := &Result{}
	for i := 1; i <= 3; i++ {
		c := curl.NewCurl(cfg.Server)
		c.PostBytes = post
		resp, err := c.Request(context.Background())
		if err != nil {
			continue
		}

		if resp.Headers["Status-Code"] != "200" {
			continue
		}

		err = json.Unmarshal(resp.Body, result)
		if err != nil {
			continue
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

	req.R.M = "upload"
	req.R.P = w

	post, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json marshal error")
	}

	result := &Result{}
	for i := 1; i <= 3; i++ {
		c := curl.NewCurl(cfg.Server)
		c.PostBytes = post
		resp, err := c.Request(context.Background())
		if err != nil {
			continue
		}

		if resp.Headers["Status-Code"] != "200" {
			continue
		}

		err = json.Unmarshal(resp.Body, result)
		if err != nil {
			continue
		}

		if result.Response.ErrMsg != "" {
			return nil, fmt.Errorf("read error: %s", result.Response.ErrMsg)
		}

		return &result.Response, nil
	}
	return nil, fmt.Errorf("upload error")
}

type DelRequest struct {
	GroupName string `json:"groupname"`
	FileName  string `json:"filename"`
}

func Del(d ReadRequest) (bool, error) {
	if d.GroupName == "" || d.FileName == "" {
		return false, fmt.Errorf("params error")
	}

	req.R.M = "del"
	req.R.P = d

	post, err := json.Marshal(req)
	if err != nil {
		return false, fmt.Errorf("json marshal error")
	}

	result := &Result{}
	for i := 1; i <= 3; i++ {
		c := curl.NewCurl(cfg.Server)
		c.PostBytes = post
		resp, err := c.Request(context.Background())
		if err != nil {
			continue
		}

		if resp.Headers["Status-Code"] != "200" {
			continue
		}

		err = json.Unmarshal(resp.Body, result)
		if err != nil {
			continue
		}

		if result.Response.ErrMsg != "" {
			return false, fmt.Errorf("read error: %s", result.Response.ErrMsg)
		}

		return true, nil
	}

	return false, fmt.Errorf("del error")
}
