package main

import (
	"context"
	"fmt"

	"github.com/gosom/kit/core"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/es/kafka"
	"github.com/gosom/kit/examples/todo"
)

func main() {
	registry := es.NewRegistry()
	todo.Register(registry)

	kafkaCfg := kafka.KafkaConfig{
		Servers:  "localhost:9092",
		GroupID:  "todo",
		Producer: true,
	}
	cfgMap := kafka.NewKafkaConfigMap(kafkaCfg)
	dispatcher, err := kafka.NewDispatcher(
		cfgMap,
		false,
		todo.COMMAND_TOPIC,
		todo.DOMAIN,
		registry,
	)
	if err != nil {
		panic(err)
	}
	defer dispatcher.Close()

	for i := 0; i < 10000; i++ {
		newTodoCmd := todo.CreateTodo{
			ID:    core.NewUUID(),
			Title: fmt.Sprintf("My %d todo", i),
		}
		ctx := context.Background()
		if cmdID, err := dispatcher.DispatchCommand(ctx, &newTodoCmd); err != nil {
			panic(err)
		} else {
			fmt.Println("Command ID: ", cmdID)
		}
	}
}
