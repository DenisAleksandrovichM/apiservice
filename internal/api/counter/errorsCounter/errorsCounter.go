package errorsCounter

import (
	"expvar"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

var c *errorsCounter

type errorsCounter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *errorsCounter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *errorsCounter) IncReq() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *errorsCounter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func Inc() {
	c.Inc()
}

func init() {
	c = &errorsCounter{m: &sync.RWMutex{}}
	expvar.Publish("Errors counter", c)
}
