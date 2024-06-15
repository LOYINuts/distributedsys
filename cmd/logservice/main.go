package main

import (
	"context"
	"distributedsys/mylog"
	"distributedsys/registry"
	"distributedsys/service"
	"fmt"
	"log"
)

func main() {
	mylog.Run("./distributed.log")
	host, port := "localhost", "8888"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		mylog.RegisterHandlers,
	)
	// 服务没有注册成功的话会返回错误服务会进行关闭
	if err != nil {
		log.Fatalln(err)
	}
	// Done返回一个通道，然后用左箭头取走通道里面的东西，
	// 如果cancel没执行则返回的通道里面没有东西，则取的操作会阻塞
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
