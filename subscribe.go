package songoku

import (
	"bytes"
	//"fmt"
)

func (c *Client) Subscribe(topic string, pid int, qos int) error {
	// MQTT Control Packet type, 0x82 mean "SUBSCRIBE"
	// -----------------------------------------
	// | Bit   | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
	// +---------------------------------------+
	// | byte1 | 1 | 0 | 0 | 0 | 0 | 0 | 1 | 0 |
	// +---------------------------------------+
	// |   0x  |       8       |       2       |
	// -----------------------------------------

	fixedheader := bytes.NewBuffer([]byte{0x82})

	// ------------------------------------------------------------
	// start encode Variable header, the version:3.1.1 MQTT
	variableheader := bytes.NewBuffer([]byte{})

	// Packet Identifier
	variableheader.WriteByte(byte((pid >> 8) & 0xFF))
	variableheader.WriteByte(byte((pid >> 0) & 0xFF))

	// topic filter
	l := len(topic)
	variableheader.WriteByte(byte((l >> 8) & 0xFF))
	variableheader.WriteByte(byte((l >> 0) & 0xFF))
	variableheader.WriteString(topic)

	// qos
	variableheader.WriteByte(byte(1 << uint(qos)))

	ll := encodeLen(variableheader.Len())
	fixedheader.Write(ll)
	fixedheader.Write(variableheader.Bytes())
	err := c.Write(fixedheader.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Unsubscribe(topic string, pid int) error {
	// MQTT Control Packet type, 0x82 mean "UNSUBSCRIBE"
	// -----------------------------------------
	// | Bit   | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
	// +---------------------------------------+
	// | byte1 | 1 | 0 | 1 | 0 | 0 | 0 | 1 | 0 |
	// +---------------------------------------+
	// |   0x  |       a       |       2       |
	// -----------------------------------------

	fixedheader := bytes.NewBuffer([]byte{0xa2})

	// ------------------------------------------------------------
	// start encode Variable header, the version:3.1.1 MQTT
	variableheader := bytes.NewBuffer([]byte{})

	// Packet Identifier
	variableheader.WriteByte(byte((pid >> 8) & 0xFF))
	variableheader.WriteByte(byte((pid >> 0) & 0xFF))

	// topic want to unsubscribe
	l := len(topic)
	variableheader.WriteByte(byte((l >> 8) & 0xFF))
	variableheader.WriteByte(byte((l >> 0) & 0xFF))
	variableheader.WriteString(topic)

	ll := encodeLen(variableheader.Len())
	fixedheader.Write(ll)
	fixedheader.Write(variableheader.Bytes())
	err := c.Write(fixedheader.Bytes())
	if err != nil {
		return err
	}

	return nil
}
