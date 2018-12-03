package dfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/logger"
	"net/http"
)

var cfg *config.DfsConfig

func init() {
	cfg = config.GetConfig().Dfs
	if cfg == nil {
		panic("Dfs config is nil")
	}
}

type Dfs struct {
	GroupName   string `json:"groupname"`
	FileName    string `json:"filename"`
	OffSet      int    `json:"offset"`
	Length      int    `json:"length"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Ext 		string `json:"ext"`
}

func (d *Dfs) Read() (string, error) {
	if d.GroupName == "" || d.FileName == "" {
		logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "groupname or filename is empty")
		return "", fmt.Errorf("groupname or filename is empty")
	}

	reqBody := commonHTTP.MakeReq(nil)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "download"
	reqBody.R.P = d

	post, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("params %s; error:%s", fmt.Sprintf("%#v", d), "json marshal error")
		return "", fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		if response.StatusCode != 200 {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "status not 200")
			continue
		}

		body := ""
		err = commonHTTP.GetRsp(response.Body, &body)
		if err != nil {
			logger.Errorf("params %s; error:%s", fmt.Sprintf("%#v", d), err.Error())
			return "", err
		}

		return body, nil
	}

	logger.Errorf("params %s; error:%s", fmt.Sprintf("%#v", d), "read error")
	return "", fmt.Errorf("read error")
}

func (d *Dfs) Write() (*Dfs, error) {
	if d.Content == "" {
		logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "content is empty")
		return nil, fmt.Errorf("params error")
	}

	reqBody := commonHTTP.MakeReq(nil)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "upload"
	reqBody.R.P = d

	post, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("params %s; error:%s", fmt.Sprintf("%#v", d), "json marshal error")
		return nil, fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		if response.StatusCode != 200 {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "status not 200")
			continue
		}

		err = commonHTTP.GetRsp(response.Body, d)
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			return nil, err
		}

		if d.GroupName == "" || d.FileName == "" {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "upload error")
			continue
		}

		return d, nil
	}

	logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "upload error")
	return nil, fmt.Errorf("upload error")
}

func (d *Dfs) Del() (bool, error) {
	if d.GroupName == "" || d.FileName == "" {
		logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "groupname or filename is empty")
		return false, fmt.Errorf("params error")
	}

	reqBody := commonHTTP.MakeReq(nil)

	reqBody.R.C = "Dfs"
	reqBody.R.M = "del"
	reqBody.R.P = d

	post, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "json marshal error")
		return false, fmt.Errorf("json marshal error")
	}

	for i := 1; i <= 3; i++ {
		request, err := http.NewRequest("POST", cfg.Server, bytes.NewBuffer(post))
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			continue
		}
		if response.StatusCode != 200 {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "status not 200")
			continue
		}

		err = commonHTTP.GetRsp(response.Body, nil)
		if err != nil {
			logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), err.Error())
			return false, err
		}

		return true, nil
	}

	logger.Errorf("params %s error:%s", fmt.Sprintf("%#v", d), "del error")
	return false, fmt.Errorf("del error")
}