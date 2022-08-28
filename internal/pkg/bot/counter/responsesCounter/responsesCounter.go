package responsesCounter

import (
	"expvar"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

var c *responsesCounter

type responsesCounter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *responsesCounter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *responsesCounter) IncReq() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *responsesCounter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func Inc() {
	c.Inc()
}

func init() {
	c = &responsesCounter{m: &sync.RWMutex{}}
	expvar.Publish("responses counter", c)
}
