package mqtt

import (
	"errors"
)

type MqttMsg struct {
	Topic   string
	Content string
	Pid     int
	Dup     bool
	Qos     int
	Retain  bool
}

func encodeLen(x int) []byte {
	l := make([]byte, 0)
	var e int
	for x > 0 {
		e = x % 128
		x /= 128
		if x > 0 {
			e = e | 128
		}
		l = append(l, byte(e))
	}
	return l
}

func deCodeLen(l []byte) (int, int, error) {
	if len(l) == 0 {
		return -1, -1, errors.New("data empty")
	}
	x := 0
	m := 1
	i := 0
	for {
		x += (int(l[i]) & 127) * m
		m *= 128
		if m > 128*128*128 {
			return -1, -1, errors.New("The size cross border")

		}
		if l[i]&128 != 128 {
			break
		}
		i++
	}
	return x, i + 1, nil
}
