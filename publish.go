package songoku

import (
	"bytes"
	"errors"
)

// for more detail of parames,
// see "http://blog.mcxiaoke.com/mqtt/protocol/MQTT-3.1.1-CN.html#pfc"
// or "http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html"
func (c *Client) Publish(msg *MqttMsg) error {
	// MQTT Publish fixed header
	// -----------------------------------------
	// | Bit   | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
	// +---------------------------------------+
	// | byte1 |               |dup|  qos  |retain|
	// +---------------------------------------+
	// |       | 0 | 0 | 1 | 1 | x | x | x | x |
	// +---------------------------------------+
	// |   0x  |       3       |       x       |
	// -----------------------------------------
	// x mean can edit

	if msg.Qos == 0 && msg.Dup {
		return errors.New("MUST set DUP=0 when qos =0")
	}

	t := 0x30
	if msg.Dup {
		t += 1 << 3
	}
	if msg.Qos == 1 {
		t += 1 << 1
	} else if msg.Qos == 2 {
		t += 1 << 2
	}
	if msg.Retain {
		t += 1
	}

	fixedheader := bytes.NewBuffer([]byte{byte(t)})

	// ------------------------------------------------------
	// start encode Variable header, the version:3.1.1 MQTT
	variableheader := bytes.NewBuffer([]byte{})

	// encode topic first
	l := len(msg.Topic)
	variableheader.WriteByte(byte((l >> 8) & 0xFF))
	variableheader.WriteByte(byte((l >> 0) & 0xFF))
	variableheader.WriteString(msg.Topic)

	// encode Packet Identifier if need
	// if qos is 0, then you can set pid anynum you want
	if msg.Qos > 0 {
		variableheader.WriteByte(byte((msg.Pid >> 8) & 0xFF))
		variableheader.WriteByte(byte((msg.Pid >> 0) & 0xFF))
	}

	// encode pyload
	variableheader.WriteString(msg.Content)

	ll := encodeLen(len(variableheader.Bytes()))
	fixedheader.Write(ll)
	fixedheader.Write(variableheader.Bytes())

	if err := c.Write(fixedheader.Bytes()); err != nil {
		return err
	}
	return nil
}
