package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/admin"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/signature"
	"io"
	"net/http"
	"time"
)

var (
	cfg *config.URIConfig
)

func init() {
	cfg = config.GetConfig().URI
	if cfg == nil {
		panic("botKit-products products config is nil")
	}
}

func getURI(productID int) (uri string, err error) {
	switch productID {
	case 1:
		uri = cfg.ChatBot
	case 3:
		uri = cfg.WinMode
	case 4:
		uri = cfg.Tsketch
	default:
		err = fmt.Errorf("botKit products error productID:%d is not defind",
			productID)
	}
	return
}

type Product struct {
	ID int `json:"product_id"`
}

func GetProduct(productID int) (*Product, error) {
	return &Product{
		ID: productID,
	}, nil
}

func GetProductExpire(managerID int, productID int) (expire time.Time, err error) {
	body, err := json.Marshal(commonHTTP.MakeReq(&struct {
		ManagerID int `json:"company_id"`
	}{managerID}))
	if err != nil {
		err = fmt.Errorf("GetProductExpire json marshal error %v", err)
		return
	}
	req, err := admin.AdminPOST("/companies/getproducts", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("GetProductExpire try make A-node request error %v", err)
		return
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("GetProductExpire A-node request error %v", err)
		return
	}
	reply := []struct {
		IsDeleted int `json:"is_deleted"`
		Deadlime  int `json:"deadline"`
		ProductID int `json:"product_id"`
	}{}
	err = commonHTTP.GetRsp(rsp.Body, &reply)
	if err != nil {
		err = fmt.Errorf("GetProductExpire A-node response error %v", err)
		return
	}

	deadTime := -1
	for _, ro := range reply {
		if ro.IsDeleted != 0 {
			continue
		}
		if ro.ProductID != productID {
			continue
		}
		deadTime = ro.Deadlime
	}
	if deadTime < 0 {
		err = fmt.Errorf("product is not exist in A-node managerID:%v productID:%v",
			managerID, productID)
		return
	}
	if deadTime == 0 {
		// no limit
		expire = time.Now().Add(time.Duration(1) * time.Hour)
		return
	}

	expire = time.Unix(int64(deadTime), 0)
	return
}

func ProductPOST(productID int, subURI string, body io.Reader) (*http.Request, error) {
	basicURI, err := getURI(productID)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", basicURI+subURI, body)
	if err != nil {
		return nil, err
	}
	err = signature.AddSignature(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
