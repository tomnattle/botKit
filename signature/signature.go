package signature

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/util"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

var cfg *config.SignatureConfig

func init() {
	cfg = config.GetConfig().Signature
	if cfg == nil {
		panic("signature config error")
	}
}

func AddSignature(req *http.Request) error {
	if req == nil {
		return fmt.Errorf("botKit signature nil http-request")
	}
	timeStamp := time.Now().Format("2006010215")
	req.Header.Add("timeStamp", timeStamp)
	nonce := util.RandStr(15)
	req.Header.Add("nonce", nonce)
	signature, err := Signature(timeStamp, nonce, cfg.SecretKey)
	if err != nil {
		return err
	}
	req.Header.Add("signature", signature)
	return nil
}

func Signature(timeStamp string, nonce string, sourcekey string) (string, error) {
	signatureSource := []string{sourcekey, timeStamp, nonce}
	if timeStamp == "" || nonce == "" {
		return "", errors.New("timeStamp or nonce can not be null")
	}
	// 时间戳 有效期限制 6个小时
	const longForm = "2006010215"
	serverTimeStamp := time.Now().Format("2006010215")
	serverTime, err := time.Parse(longForm, serverTimeStamp)
	if err != nil {
		return "", fmt.Errorf("Error in time format %v should be 2006010215", err)
	}
	timeStr, err := time.Parse(longForm, timeStamp)
	if err != nil {
		return "", fmt.Errorf("Error in time format %v should be 2006010215", err)
	}
	if (timeStr.Sub(serverTime) > 6*time.Hour) || (serverTime.Sub(timeStr) > 6*time.Hour) {
		return "", errors.New("time is expires")
	}

	sort.Slice(signatureSource, func(i, j int) bool { return signatureSource[i] < signatureSource[j] })
	sha1er := sha1.New()
	io.WriteString(sha1er, strings.Join(signatureSource, ""))
	signatureStr := fmt.Sprintf("%x", sha1er.Sum(nil))
	return signatureStr, nil
}

func VerifySignature(req *http.Request) (pass bool, err error) {
	timeStamp := req.Header.Get("timeStamp")
	signatureStr := req.Header.Get("signature")
	nonce := req.Header.Get("nonce")
	if signatureStr == "" {
		err = fmt.Errorf("signature can not be null")
		return
	}
	selfSignatureStr, err := Signature(timeStamp, nonce, cfg.SecretKey)
	if err != nil {
		err = fmt.Errorf("signature create err %v", err)
		return
	}
	if selfSignatureStr == signatureStr {
		pass = true
		return
	}
	return
}
