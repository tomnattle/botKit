package products

import (
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/Redis"
	"github.com/ifchange/botKit/logger"
)

const (
	ttl = 30
)

func key(managerID int) string {
	return Redis.FormatKey(fmt.Sprintf("manager-products-stat e-cache-%d", managerID))
}

func getProductStateCache(conn *Redis.RedisCommon, id int) ([]*ProductState, bool) {
	data, err := conn.Cmd("GET", key(id)).Bytes()
	if err != nil {
		logger.Printf("manager products try exec redis query error %v", err)
		return nil, false
	}

	pss := make([]*ProductState, 0)
	err = json.Unmarshal(data, &pss)
	if err != nil {
		logger.Printf("manager products try exec redis query error %v", err)
		return nil, false
	}
	return pss, true
}

func saveProductStateCache(conn *Redis.RedisCommon, id int, pss []*ProductState) {
	data, err := json.Marshal(pss)
	if err != nil {
		return
	}

	err = conn.Cmd("SETEX", key(id), ttl, data).Err
	if err != nil {
		logger.Printf("manager products try exec redis save error %v", err)
		return
	}
}
