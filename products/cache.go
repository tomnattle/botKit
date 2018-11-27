package products

import (
	"fmt"
	"github.com/ifchange/botKit/cache"
	"time"
)

var (
	productStateCache *cache.Cache
	cacheTTL          = time.Duration(30) * time.Second
)

func init() {
	var err error
	productStateCache, err = cache.New()
	if err != nil {
		panic(err)
	}
}

func key(managerID int) string {
	return fmt.Sprintf("manager-products-state-cache-%d", managerID)
}

func getProductStateCache(id int) ([]*ProductState, bool) {
	pps, ok := productStateCache.Get(key(id)).([]*ProductState)
	if !ok {
		return nil, false
	}
	return pps, true
}

func saveProductStateCache(id int, pss []*ProductState) {
	productStateCache.Set(key(id), pss, cacheTTL)
}
