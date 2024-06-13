package main

import (
	"context"
	"distributedsys/mylog"
	"distributedsys/service"
	"fmt"
	"log"
)

func main() {
	mylog.Run("./distributed.log")
	host, port := "localhost", "8888"
	ctx, err := service.Start(
		context.Background(),
		"Log Service",
		host,
		port,
		mylog.RegisterHandlers,
	)
	if err != nil {
		log.Fatalln(err)
	}
	// Done返回一个通道，然后用左箭头取走通道里面的东西，
	// 如果cancel没执行则返回的通道里面没有东西，则取的操作会阻塞
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
