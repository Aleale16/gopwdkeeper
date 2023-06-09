package grpcserver

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "pwdkeeper/internal/app/proto"
	"pwdkeeper/internal/app/storage"

	"google.golang.org/grpc"
)

// Grpcserverstart starts gRPC server
func Grpcserverstart() {
	// func Grpcserverstart() (error) {

	storage.Initdb()

	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	S := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterActionsServer(S, &ActionsServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	//	if err := s.Serve(listen); err != nil {
	//		log.Fatal(err)
	//	}
	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// Ожидаем события от ОС
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// сообщаем об ошибках в канал
	go func() {
		if err := S.Serve(listen); err != nil {
			errChan <- err
		}
	}()
	defer func() {
		S.GracefulStop()
	}()

	select {
	case err := <-errChan:
		log.Printf("Fatal error: %v\n", err)
	case <-stopChan:
	}
}
