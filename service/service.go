package service

import (
	"context"
	"distributedsys/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName,
	host, port string) context.Context {
	ctxtmp, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = host + ":" + port
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctxtmp)
		cancel()
	}()
	return ctxtmp
}
