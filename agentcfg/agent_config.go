package agentcfg

import (
	"strings"
	"sync"
	"time"
)

var PollInterval = time.Minute

type syncMap struct {
	sync.RWMutex
	m map[string]map[string]interface{}
}

func (sm *syncMap) get(k string) map[string]interface{} {
	sm.RLock()
	defer sm.RUnlock()
	ret, _ := sm.m[k]
	return ret
}

func (sm *syncMap) mset(m map[string]map[string]interface{}) {
	sm.Lock()
	defer sm.Unlock()
	for k, v := range m {
		sm.m[k] = v
	}
}

// takes a list of nested keys and returns the associated value and an UID
type LookupFunc func(...string) (map[string]interface{}, string)

var once sync.Once

func StartPolling(remote RemoteSource) func(...string) (map[string]interface{}, string) {

	var cfgMap *syncMap
	var lookup LookupFunc

	once.Do(func() {
		cfgMap = &syncMap{}
		lookup = func(keys ...string) (map[string]interface{}, string) {
			ret := cfgMap.get(concat(keys...))
			id, ok := ret["id"]
			if !ok {
				return ret, ""
			}
			etag, ok := id.(string)
			if !ok {
				return ret, ""
			}
			return ret, etag
		}
	})

	ticker := time.NewTicker(PollInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cfgMap.mset(remote.Fetch())
			}
		}
	}()
	return lookup
}

func concat(args ...string) string {
	return strings.Join(args, ":")
}
