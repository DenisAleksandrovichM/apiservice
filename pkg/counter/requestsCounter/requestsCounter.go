package requestsCounter

import (
	"expvar"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

var c *requestsCounter

type requestsCounter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *requestsCounter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *requestsCounter) IncReq() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *requestsCounter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func Inc() {
	c.Inc()
}

func init() {
	c = &requestsCounter{m: &sync.RWMutex{}}
	expvar.Publish("Requests counter", c)
}
