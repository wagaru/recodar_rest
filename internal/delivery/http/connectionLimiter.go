package http

type ConnectionLimiter struct {
	concurrent int
	pool       chan bool
}

func newConnectionLimiter(concurrent int) *ConnectionLimiter {
	return &ConnectionLimiter{
		concurrent: concurrent,
		pool:       make(chan bool, concurrent),
	}
}

func (c *ConnectionLimiter) GetConnection() bool {
	if len(c.pool) >= c.concurrent {
		return false
	}
	c.pool <- true
	return true
}

func (c *ConnectionLimiter) ReleaseConnection() {
	<-c.pool
}
