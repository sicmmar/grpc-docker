package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"bytes"
	"os"
	"net/http"
	"encoding/json"

	"google.golang.org/grpc"
	"grpcserver/user.pb/user.pb"
)

type servidor struct{}

type usuario struct{
	Name string	`json:"name"`
	Location string `json:"location"`
	Age int `json:"age"`
	Infectedtype string `json:"infectedtype"`
	State string `json:"state"`
	Way string `json:"way"`
}

func (*servidor) RegUser(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserResponse, error) {
	fmt.Println("Todo bien!")

	cuerpoPeticion, _ := json.Marshal(usuario{
		Name: req.User.Name,
		Location: req.User.Location,
		Age: int(req.User.Age),
		Infectedtype: req.User.Infectedtype,
		State: req.User.State,
		Way: "GRPC",
	})

	pet := bytes.NewBuffer(cuerpoPeticion)

	resp, err := http.Post("http://35.222.55.115:8080/nuevoRegistro", "application/json", pet)
	if err != nil {
		log.Fatalln("Error al registrar nuevo: ", err)
	}

	defer resp.Body.Close()

	cuerpo, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Fatalln(err)
	}

	result := &user_pb.UserResponse{
		Resultado: string(cuerpo),
	}

	return result, nil
}

func main() {
	host := os.Getenv("HOST")
	//host := "0.0.0.0:8081"
	fmt.Println("Servidor iniciado en ", host)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("F con el servidor: %v", err)
	}

	fmt.Println("Empezando servidor grpc ...")

	s := grpc.NewServer()

	user_pb.RegisterUserServiceServer(s, &servidor{})

	fmt.Println("Servidor a la espera ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("El servidor no funciona: %v", err)
	}
}
