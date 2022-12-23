package connection

import "sync"

type Client interface {
	Connection() Connection
	SetConnection(conn Connection)
	Online() bool
}

type client struct {
	currentConn Connection
	connMu sync.Mutex
	isOnline bool
}

func NewClient() *client {
	c := new(client)
	c.isOnline = false
	c.currentConn = NewNopConnection()
	return c
}

func (c *client) Connection() Connection {
	return c.currentConn
}

func (c *client) SetConnection(conn Connection) {
	c.connMu.Lock()
	defer c.connMu.Unlock()

	if conn == nil {
		c.isOnline = false
		c.currentConn = NewNopConnection()
	} else {
		c.isOnline = true
		c.currentConn = conn

		go func() {
			<-conn.Context().Done()
			c.SetConnection(nil)
		}()
	}
}

func (c *client) Online() bool {
	return c.isOnline
}
