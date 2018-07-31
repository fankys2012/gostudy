package orm

import "sync"

var (
	modelCache = &_modelCache{
		cache:           make(map[string]*modelInfo),
		cacheByFullName: make(map[string]*modelInfo),
	}
)

type _modelCache struct {
	sync.RWMutex
	orders          []string
	cache           map[string]*modelInfo
	cacheByFullName map[string]*modelInfo
	done            bool
}

func (this *_modelCache) set(table string, mi *modelInfo) *modelInfo {
	mii := this.cache[table]
	this.cache[table] = mi
	this.cacheByFullName[mi.fullName] = mi
	if mii == nil {
		this.orders = append(this.orders, table)
	}
	return mii
}

func (this *_modelCache) get(table string) (*modelInfo, bool) {
	mi, ok := this.cache[table]
	return mi, ok
}
