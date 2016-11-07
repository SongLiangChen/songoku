package mqtt

import (
	"net"
	"sync"
)

// the mqtt client, base on mqtt version 3.1.1
// for more detail, please see "http://blog.mcxiaoke.com/mqtt/protocol/MQTT-3.1.1-CN.html#pfc" or "http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html"

type onpublishfunc func(*MqttMsg)
type onconnackfunc func()
type onpubackfunc func()
type onpubrecfunc func(pid []byte)
type onpubrelfunc func(pid []byte)
type onpubcompfunc func(pid []byte)
type onsubackfunc func()
type onunsubackfunc func()
type onpingrespfunc func()

type Client struct {
	ID           string // the Client Identifier, cannot be nil
	Keepalive    int    // idle time, when timeout, disconnect from server
	WillRetain   bool   // "will retain" flag, it must be false unless "WillFlag == true"
	WillQos      int    // "will qos" flag, it must be zero unless "WillFlag == frue"
	WillFlag     bool   // "will flag"
	CleanSession bool

	WillTopic  string
	WillWorlds string
	UserName   string
	PassWorld  string

	conn *net.TCPConn

	// expand this function, and assign to client
	onPublish  onpublishfunc
	onConnack  onconnackfunc
	onPubAck   onpubackfunc
	onPubRec   onpubrecfunc
	onPubRel   onpubrelfunc
	onPubComp  onpubcompfunc
	onSubAck   onsubackfunc
	onUnsubAck onunsubackfunc
	onPingResp onpingrespfunc

	sync.Mutex
}

func NewMqttClient() *Client {
	return &Client{
		ID:           "default_id",
		Keepalive:    60,
		WillRetain:   false,
		WillQos:      0,
		WillFlag:     false,
		CleanSession: true,

		WillTopic:  "",
		WillWorlds: "",
		UserName:   "",
		PassWorld:  "",

		onPublish:  func(*MqttMsg) {},
		onConnack:  func() {},
		onPubAck:   func() {},
		onPubRec:   func(pid []byte) {},
		onPubRel:   func(pid []byte) {},
		onPubComp:  func(pid []byte) {},
		onSubAck:   func() {},
		onUnsubAck: func() {},
		onPingResp: func() {},
	}
}

func (c *Client) HandPublish(fn func(msg *MqttMsg)) {
	c.onPublish = fn
}

func (c *Client) HandConnack(fn func()) {
	c.onConnack = fn
}

func (c *Client) HadPuback(fn func()) {
	c.onPubAck = fn
}

func (c *Client) HandPubrec(fn func(pid []byte)) {
	c.onPubRec = fn
}

func (c *Client) HandPubrel(fn func(pid []byte)) {
	c.onPubRel = fn
}

func (c *Client) HandPubcomp(fn func(pid []byte)) {
	c.onPubComp = fn
}

func (c *Client) HandSuback(fn func()) {
	c.onSubAck = fn
}

func (c *Client) Write(d []byte) error {
	c.Lock()
	defer c.Unlock()
	_, err := c.conn.Write(d)
	return err
}
