package songoku

import (
	"bytes"
	"errors"
	"net"
)

func (c *Client) Connect(address string) error {
	if c.ID == "" {
		return errors.New("mqtt client id can not be nil")
	}

	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return err
	}

	if !c.WillFlag {
		if c.WillRetain == true || c.WillQos > 0 {
			return errors.New("WillFlag is 0, you should make WillRetain and WillQos 0")
		}
	} else {
		if c.WillTopic == "" || c.WillWorlds == "" {
			return errors.New("WillFlag is 1, you must set WillTopic and WillWorlds")
		}
	}

	if c.UserName == "" && c.PassWorld != "" {
		return errors.New("no username but has passworld")
	}

	// MQTT Control Packet type, 0x10 mean "CONNECT"
	// -----------------------------------------
	// | Bit   | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
	// +---------------------------------------+
	// | byte1 | 0 | 0 | 0 | 1 | 0 | 0 | 0 | 0 |
	// +---------------------------------------+
	// |   0x  |       1       |       0       |
	// -----------------------------------------
	fixedheader := bytes.NewBuffer([]byte{0x10})

	// ------------------------------------------------------------
	// start encode Variable header, the version:3.1.1 MQTT
	variableheader := bytes.NewBuffer([]byte{0x00, 0x04, 0x4d, 0x51, 0x54, 0x54, 0x04})

	// encode Connect Flags
	cf := 0
	if c.CleanSession {
		cf += 2
	}
	if c.WillFlag {
		cf += 4
	}
	if c.WillQos == 1 {
		cf += 8
	} else if c.WillQos == 2 {
		cf += 16
	}
	if c.WillRetain {
		cf += 32
	}
	if c.PassWorld != "" {
		cf += 64
	}
	if c.UserName != "" {
		cf += 128
	}
	variableheader.WriteByte(byte(cf))

	// encode keepalive
	variableheader.WriteByte(byte((c.Keepalive >> 8) & 0xFF))
	variableheader.WriteByte(byte((c.Keepalive >> 0) & 0xFF))

	// the entire Variable header encode finish
	// -------------------------------------------------------------

	// -------------------------------------------------------------
	// start encode Payload
	payload := bytes.NewBuffer([]byte{})
	// encode Client Identifier
	l := len(c.ID)
	payload.WriteByte(byte((l >> 8) & 0xFF))
	payload.WriteByte(byte((l >> 0) & 0xFF))
	payload.WriteString(c.ID)

	if c.WillFlag {
		l = len(c.WillTopic)
		payload.WriteByte(byte((l >> 8) & 0xFF))
		payload.WriteByte(byte((l >> 0) & 0xFF))
		payload.WriteString(c.WillTopic)

		l = len(c.WillWorlds)
		payload.WriteByte(byte((l >> 8) & 0xFF))
		payload.WriteByte(byte((l >> 0) & 0xFF))
		payload.WriteString(c.WillWorlds)
	}

	if c.UserName != "" {
		l = len(c.UserName)
		payload.WriteByte(byte((l >> 8) & 0xFF))
		payload.WriteByte(byte((l >> 0) & 0xFF))
		payload.WriteString(c.UserName)
	}

	if c.PassWorld != "" {
		l = len(c.PassWorld)
		payload.WriteByte(byte((l >> 8) & 0xFF))
		payload.WriteByte(byte((l >> 0) & 0xFF))
		payload.WriteString(c.PassWorld)
	}

	// the entire Payload encode finish
	// -------------------------------------------------------------

	l = variableheader.Len() + payload.Len()
	ll := encodeLen(l)
	fixedheader.Write(ll)
	fixedheader.Write(variableheader.Bytes())
	fixedheader.Write(payload.Bytes())

	err = c.Write(fixedheader.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Disconnect() error {
	// MQTT Control Packet type, 0xe0 mean "DISCONNECT"
	// -----------------------------------------
	// | Bit   | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
	// +---------------------------------------+
	// | byte1 | 1 | 1 | 1 | 0 | 0 | 0 | 0 | 0 |
	// +---------------------------------------+
	// |   0x  |       e       |       0       |
	// -----------------------------------------
	return c.Write([]byte{0xe0, 0x00})
}
