# go-amqp-tasks

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rafaeltokyo/go-amqp-tasks/blob/master/LICENSE)



## Why?

I need to simple tasks with amqp for my projects


## go-amqp-tasks woker

```go
package main

import (
	"fmt"

	goamqptasks "github.com/rafaeltokyo/go-amqp-tasks"
)

func handleHello(action string, body []byte) {
	fmt.Printf("Body %s", body)
}

func main() {

	tasks, err := goamqptasks.NewTask("localhost:5672")
	if err != nil {
		panic(err)
	}
	forever := make(chan bool)
	go tasks.On("handleHello", handleHello)
	<-forever

}
```


## go-amqp-tasks publish

```go
package main

import (
	"fmt"

	goamqptasks "github.com/rafaeltokyo/go-amqp-tasks"
)

func main() {

	tasks, err := goamqptasks.NewTask("localhost:5672")
	if err != nil {
		panic(err)
	}
	tasks.Publish("handleHello", "hello")

}
}
```
