package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/Redis"
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

func getProductInfo(productID int) (name string, desc string, sort int, err error) {
	body, err := json.Marshal(commonHTTP.MakeReq(&struct {
		ProductID int `json:"id"`
	}{productID}))
	if err != nil {
		err = fmt.Errorf("getProductInfo json marshal error %v", err)
		return
	}
	req, err := admin.AdminPOST("/products/detail", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("getProductInfo try make A-node request error %v", err)
		return
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("getProductInfo A-node request error %v", err)
		return
	}
	reply := struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
		Sort int    `json:"sort"`
	}{}
	err = commonHTTP.GetRsp(rsp.Body, &reply)
	if err != nil {
		err = fmt.Errorf("getProductInfo A-node response error %v", err)
		return
	}
	name = reply.Name
	desc = reply.Desc
	sort = reply.Sort
	return
}

type Product struct {
	ID        int       `json:"product_id"`
	Name      string    `json:"product_name"`
	Desc      string    `json:"desc"`
	Sort      int       `json:"sort"`
	IsDeleted int       `json:"is_deleted"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func GetProduct(productID int) (*Product, error) {
	name, desc, sort, err := getProductInfo(productID)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:   productID,
		Name: name,
		Desc: desc,
		Sort: sort,
	}, nil
}

func GetProductExpire(managerID int, productID int) (expire time.Time, err error) {
	reply, err := getProductStateByManagerID(managerID)
	if err != nil {
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
		deadTime = ro.Deadline
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

func GetProductsFromAdmin() ([]*Product, error) {
	return getProductsFromAdmin()
}

func getProductsFromAdmin() ([]*Product, error) {
	body, err := json.Marshal(commonHTTP.MakeReq(&struct{}{}))
	if err != nil {
		err = fmt.Errorf("getProductsFromAdmin json marshal error %v", err)
		return nil, err
	}
	req, err := admin.AdminPOST("/products/search", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("getProductsFromAdmin try make A-node request error %v", err)
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("getProductsFromAdmin A-node request error %v", err)
		return nil, err
	}
	var reply = &struct {
		Total    int        `json:"total"`
		Products []*Product `json:"results"`
	}{}
	err = commonHTTP.GetRsp(rsp.Body, &reply)
	if err != nil {
		err = fmt.Errorf("getProductsFromAdmin unmarshal error %v", err)
		return nil, err
	}

	return reply.Products, nil
}

type ProductState struct {
	ProductID     int `json:"product_id"`
	IsDeleted     int `json:"is_deleted"`
	Deadline      int `json:"deadline"`
	PurchaseState int `json:"purchase_state"`
}

func (ps *ProductState) Marshal() ([]byte, error) {
	if ps == nil {
		return nil, fmt.Errorf("ProductState is nil")
	}
	return json.Marshal(&ps)
}

func (ps *ProductState) Unmarshal(data []byte) error {
	if ps == nil {
		return fmt.Errorf("ProductState is nil")
	}
	return json.Unmarshal(data, &ps)
}

func (ps *ProductState) parse() {
	if ps.Deadline == 0 {
		ps.PurchaseState = 1
	}
	tm := time.Unix(int64(ps.Deadline), 0)
	if time.Now().After(tm) {
		ps.PurchaseState = 2
	}
}

func GetProductStateByManagerID(id int) ([]*ProductState, error) {
	return getProductStateByManagerID(id)
}

func getProductStateByManagerID(id int) ([]*ProductState, error) {
	conn, err := Redis.GetRedis()
	if err == nil {
		pss, exist := getProductStateCache(conn, id)
		if exist {
			return pss, nil
		}
	}
	body, err := json.Marshal(commonHTTP.MakeReq(&struct {
		CompanyID int `json:"company_id"`
	}{id}))
	if err != nil {
		err = fmt.Errorf("getProductStateByManagerID json marshal error %v", err)
		return nil, err
	}
	req, err := admin.AdminPOST("/companies/getproducts", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("getProductStateByManagerID try makPurchaseStatee A-node request error %v", err)
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("getProductStateByManagerID A-node request error %v", err)
		return nil, err
	}
	states := []*ProductState{}
	err = commonHTTP.GetRsp(rsp.Body, &states)
	if err != nil {
		err = fmt.Errorf("getProductStateByManagerID unmarshal error %v", err)
		return nil, err
	}

	for _, s := range states {
		s.parse()
	}

	saveProductStateCache(conn, id, states)

	return states, nil
}
