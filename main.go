package main

import (
	"context"
	"fmt"
	"main/application"
	"os"
	"os/signal"
)

func main() {
	application := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := application.Start(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
