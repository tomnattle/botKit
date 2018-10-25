package SMS

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"ifchange/tsketch/kit/config"
	"ifchange/tsketch/kit/util"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

var (
	cfg *config.SMSConfig
)

func init() {
	cfg = config.GetConfig().SMS
	if cfg == nil {
		panic("SMS config is nil")
	}
}

func Request(phone, content string) error {
	req, err := payload(phone, content)
	if err != nil {
		return fmt.Errorf("kit-SMS error in-queue, try make payload %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("kit-SMS error in request, %v",
			err.Error())
	}
	if err := decodeRsp(rsp); err != nil {
		return fmt.Errorf("kit-SMS error in response, %v",
			err.Error())
	}
	return nil
}

func newCurl(request SMSRequest) *CurlCommon {
	timeStamp := time.Now().Format("200601021504")
	nonce := util.RandStr(15)

	signatureSource := []string{cfg.Secret, timeStamp, nonce}
	sort.Slice(signatureSource, func(i, j int) bool { return signatureSource[i] < signatureSource[j] })

	sha1er := sha1.New()
	io.WriteString(sha1er, strings.Join(signatureSource, ""))
	signature := fmt.Sprintf("%x", sha1er.Sum(nil))

	return &CurlCommon{
		Header: CurlCommonHeader{
			AppID:     config.GetConfig().AppID,
			LogID:     config.GetEnvironment().String(),
			TimeStamp: timeStamp,
			Nonce:     nonce,
			Signature: signature,
			IP:        "8.9.10.9",
		},
		Request: request,
	}
}

type CurlCommon struct {
	Header  CurlCommonHeader `json:"header"`
	Request SMSRequest       `json:"request"`
}

func (ins *CurlCommon) Marshal() ([]byte, error) {
	return json.Marshal(ins)
}

type CurlCommonHeader struct {
	AppID     int    `json:"app_id"`
	LogID     string `json:"log_id"`
	UID       string `json:"uid"`
	UName     string `json:"uname"`
	Provider  string `json:"provider"`
	SignID    string `json:"signid"`
	Version   string `json:"version"`
	IP        string `json:"ip"`
	TimeStamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type SMSRequest struct {
	C string           `json:"c"`
	M string           `json:"m"`
	P SMSRequestInside `json:"p"`
}

type SMSRequestInside struct {
	Basic      SMSRequestBasic      `json:"basic"`
	BasicExtra SMSRequestBasicExtra `json:"basic_extra"`
}

type SMSRequestBasic struct {
	AppID      int    `json:"app_id"`
	AppType    int    `json:"app_type"`
	MsgType    string `json:"msg_type"`
	SendType   int    `json:"send_type"`
	IsCallback string `json:"is_callback"`
	ToAppID    string `json:"to_app_id"`
	ToUserID   string `json:"to_user_id"`
	ExecutedAt string `json:"executed_at"`
}

type SMSRequestBasicExtra struct {
	Content   string `json:"contents"`
	ToPhone   string `json:"to_phones"`
	Signature int    `json:"signature"`
	Account   string `json:"account"`
}

func payload(phone string, content string) (*http.Request, error) {
	data := newCurl(SMSRequest{
		C: "basic/Logic_basic",
		M: "save",
		P: SMSRequestInside{
			Basic: SMSRequestBasic{
				AppID:      config.GetConfig().AppID,
				AppType:    0,
				MsgType:    "sms",
				SendType:   0,
				IsCallback: "N",
				ToAppID:    "5",
				ToUserID:   "1",
				ExecutedAt: time.Now().Format("200601021504"),
			},
			BasicExtra: SMSRequestBasicExtra{
				Content:   content,
				ToPhone:   phone,
				Signature: 1,
				Account:   cfg.Account,
			},
		},
	})
	body, err := data.Marshal()
	if err != nil {
		return nil, fmt.Errorf("json marshal error %v, %v",
			data, err)
	}
	req, err := http.NewRequest("POST", cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

type SMSResponse struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

func decodeRsp(rsp *http.Response) error {
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("kit-SMS error read http response error %v", err)
	}
	data := struct {
		SMSResponse `json:"response"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("kit-SMS error unmarshal http response error %v", err)
	}
	if data.SMSResponse.ErrNo == 0 {
		return nil
	}
	return fmt.Errorf("kit-SMS error response %v", data.SMSResponse)
}
