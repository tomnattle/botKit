package sessionKit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/admin"
	"github.com/ifchange/botKit/commonHTTP"
	"net/http"
)

func GetSecretKey(srcID, managerID int) (string, error) {
	body := &bytes.Buffer{}
	reqBody := commonHTTP.MakeReq(&struct {
		SrcID     int `json:"src_id"`
		ManagerID int `json:"id"`
	}{
		SrcID:     srcID,
		ManagerID: managerID,
	})

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	_, err = body.Write(reqData)
	if err != nil {
		return "", fmt.Errorf("try write body error %v", err)
	}

	req, err := admin.AdminPOST("/companies/getsecretkey", body)
	if err != nil {
		return "", fmt.Errorf("admin make request error %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("admin request error %v", err)
	}
	defer rsp.Body.Close()

	secretKey := ""

	err = commonHTTP.GetRsp(rsp.Body, &secretKey)
	if err != nil {
		return "", err
	}

	if len(secretKey) == 0 {
		return "", fmt.Errorf("empty secretKey")
	}
	return secretKey, nil
}
