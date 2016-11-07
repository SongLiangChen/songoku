package mqtt

import (
	"bytes"
	log "github.com/thinkboy/log4go"
)

var (
	puback  = []byte{0x40, 0x02}
	pubrec  = []byte{0x50, 0x02}
	pubrel  = []byte{0x62, 0x02}
	pubcomp = []byte{0x70, 0x02}
)

// it usually call by go, like "go c.Consume()"
func (c *Client) Consume() {

	var n, i, j int
	var err error
	var packet_type byte

	buf := make([]byte, 65535)
	tmp := make([]byte, 2*65535)
	tmplen := 0
	for {
		if n, err = c.conn.Read(buf); err != nil {
			c.conn.Close()
			break
		}

		// tcp socket has packet splicing problem
		for i = 0; i < n; i++ {
			tmp[tmplen+i] = buf[i]
		}
		tmplen += n

		for tmplen > 1 {
			packet_type = tmp[0]
			if n, j, err = deCodeLen(tmp[1:]); err == nil {
				if tmplen-j-1 < n {
					continue
				}

				t := make([]byte, n)
				copy(t, tmp[1+j:1+j+n])

				switch packet_type >> 4 {
				// CONNACK
				case 0x02:
					go c.handleConnect(t)
					break
				// PUBLISH
				case 0x03:
					go c.handlePublish(packet_type, t)
					break
				// PUBACK
				case 0x04:
					go c.onPubAck()
					break
				// PUBREC
				case 0x05:
					go c.handlePubRec(t)
					break
				// PUBREL
				case 0x06:
					go c.handlePubRel(t)
					break
				// PUBCOMP
				case 0x07:
					go c.handlePubComp(t)
					break
				// SUBACK
				case 0x09:
					go c.handleSubAck(t)
					break
				// UNSUBACK
				case 0x0b:
					log.Info("unsubscribe topic success")
					break
				// PINGRESP
				case 0x0d:
					break
				}

				for i = n + j + 1; i < tmplen; i++ {
					tmp[i-n-j-1] = tmp[i]
				}
				tmplen -= (n + j + 1)
			}
		}
	}
}

func (c *Client) handleConnect(buf []byte) {
	// check CONNACK
	if buf[1] == 0 {
		go c.heartbeat()
		log.Info("connect mosquitto success")
	} else {
		c.conn.Close()
		log.Info("connect mosquitto fail")
		return
	}
	go c.onConnack()
}

func (c *Client) handlePublish(packet_type byte, buf []byte) {
	var pid string
	index := 0
	topiclen := int(buf[index+1]&0xFF) | int((buf[index]&0xFF)<<8)
	index += 2
	m := &MqttMsg{}
	m.Topic = string(buf[index : index+topiclen])
	index += topiclen
	if packet_type&6 > 0 {
		pid = string(buf[index : index+2])
		index += 2
		m.Content = string(buf[index:])
	} else {
		m.Content = string(buf[index:])
	}
	go c.onPublish(m)

	if packet_type&2 > 0 {
		// qos == 1
		var buffer bytes.Buffer
		buffer.Write(puback)
		buffer.WriteString(pid)
		err := c.Write(buffer.Bytes())
		if err != nil {
			c.conn.Close()
			return
		}
	} else if packet_type&4 > 0 {
		// qos == 2
		var buffer bytes.Buffer
		buffer.Write(pubrec)
		buffer.WriteString(pid)
		err := c.Write(buffer.Bytes())
		if err != nil {
			c.conn.Close()
			return
		}
	}
}

func (c *Client) handlePubRec(buf []byte) {
	pid := string(buf[:2])
	go c.onPubRec([]byte(pid))
	var buffer bytes.Buffer
	buffer.Write(pubrel)
	buffer.WriteString(pid)
	err := c.Write(buffer.Bytes())
	if err != nil {
		c.conn.Close()
		return
	}
}

func (c *Client) handlePubRel(buf []byte) {
	pid := string(buf[:2])
	go c.onPubRel([]byte(pid))
	var buffer bytes.Buffer
	buffer.Write(pubcomp)
	buffer.WriteString(pid)
	err := c.Write(buffer.Bytes())
	if err != nil {
		c.conn.Close()
		return
	}
}

func (c *Client) handlePubComp(buf []byte) {
	pid := string(buf[:2])
	go c.onPubComp([]byte(pid))
}

func (c *Client) handleSubAck(buf []byte) {
	if buf[2] == 0x80 {
		log.Info("subscribe topic fail")
	} else {
		log.Info("subscribe topic success")
	}
}
