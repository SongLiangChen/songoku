# songoku
the mqtt client by go.

demo:
```
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/SongLiangChen/songoku"
)

func main() {
	client := songoku.NewMqttClient()

	client.HandPublish(func(msg *songoku.MqttMsg) {
		fmt.Printf("rec new msg: %v\n", *msg)
	})

	if err := client.Connect("127.0.0.1:1883"); err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect()

	go client.Consume()

	client.Subscribe("topicA", 1, 0)

	time.Sleep(time.Second)
	go TestPublish()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func TestPublish() {
	client := songoku.NewMqttClient()
	client.ID = "your_client_id"
	if err := client.Connect("127.0.0.1:1883"); err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect()

	go client.Consume()

	client.Publish(&songoku.MqttMsg{
		Topic:   "topicA",
		Content: "hello",
		Qos:     0,
	})

	time.Sleep(time.Second)
}

```
