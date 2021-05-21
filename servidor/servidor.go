package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	user_pb "grpcserver/user.pb"

	"google.golang.org/grpc"
)

type servidor struct {
	user_pb.UnimplementedUserServiceServer
}

type usuario struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	VaccineType string `json:"vaccine_type"`
	Way         string `json:"way"`
}

func (*servidor) RegUser(ctx context.Context, req *user_pb.UserRequest) (*user_pb.UserResponse, error) {
	respuesta := ""

	fmt.Println("Todo bien!")

	cuerpoPeticion, _ := json.Marshal(usuario{
		Name:        req.User.Name,
		Location:    req.User.Location,
		Gender:      req.User.Gender,
		Age:         int(req.User.Age),
		VaccineType: req.User.VaccineType,
		Way:         "GRPC",
	})

	pet := bytes.NewBuffer(cuerpoPeticion)

	// ENVIAR DATOS A MONGODB
	resp, err := http.Post("http://34.66.140.170:8080/nuevoRegistro", "application/json", pet)
	if err != nil {
		log.Fatalln("Error al registrar en Mongo: ", err)
	}

	defer resp.Body.Close()

	cuerpo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	respuesta += string(cuerpo)

	// ENVIAR DATOS A REDIS
	resredis, erredis := http.Post("http://35.223.156.4:7019/nuevoRegistro", "application/json", bytes.NewBuffer(cuerpoPeticion))
	if erredis != nil {
		log.Fatalln("Error al registrar en Redis: ", erredis)
	}

	defer resredis.Body.Close()

	cRedis, eredis := ioutil.ReadAll(resredis.Body)
	if eredis != nil {
		log.Fatalln(eredis)
	}
	respuesta += string(cRedis)

	result := &user_pb.UserResponse{
		Resultado: respuesta,
		//Resultado: "PRUEBA",
	}

	return result, nil
}

func main() {
	//host := os.Getenv("HOST")
	host := "0.0.0.0:8081"
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
