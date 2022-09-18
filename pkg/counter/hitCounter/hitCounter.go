package hitCounter

import (
	"expvar"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

var c *hitCounter

type hitCounter struct {
	cnt uint
	m   *sync.RWMutex
}

func (c *hitCounter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *hitCounter) IncReq() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *hitCounter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func Inc() {
	c.Inc()
}

func init() {
	c = &hitCounter{m: &sync.RWMutex{}}
	expvar.Publish("Hit counter", c)
}
