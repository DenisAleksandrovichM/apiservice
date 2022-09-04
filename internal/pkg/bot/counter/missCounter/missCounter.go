package missCounter

import (
	"expvar"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

var c *missCounter

type missCounter struct {
	cnt uint
	m   *sync.RWMutex
}

func (c *missCounter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *missCounter) IncReq() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *missCounter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func Inc() {
	c.Inc()
}

func init() {
	c = &missCounter{m: &sync.RWMutex{}}
	expvar.Publish("Miss counter", c)
}
