package songoku

import (
	"time"
)

func (c *Client) heartbeat() {
	t := time.NewTicker(time.Second * time.Duration(c.Keepalive/2))

	var err error
	pingresq := []byte{0xc0, 0x00}

	for {
		select {
		case <-t.C:
			err = c.Write(pingresq)
			if err != nil {
				c.conn.Close()
				return
			}
		}
	}
}
