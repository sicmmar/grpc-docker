package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"

	user_pb "grpccliente/user.pb"

	"google.golang.org/grpc"
)

type userStruct struct {
	Name        string
	Location    string
	Gender      string
	Age         int64
	VaccineType string
}

func registrarUsuario(nameparam string, locationparam string, genderparam string, ageparam int64, vaccinetypeparam string) {
	//server_host := os.Getenv("SERVER_HOST")
	server_host := "0.0.0.0:8081"
	fmt.Println("Enviando peticion ...")

	cc, err := grpc.Dial(server_host, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error enviando peticion: %v", err)
	}

	defer cc.Close()

	c := user_pb.NewUserServiceClient(cc)

	fmt.Println("Todo bien en la conexion")

	request := &user_pb.UserRequest{
		User: &user_pb.Usuario{
			Name:        nameparam,
			Location:    locationparam,
			Gender:      genderparam,
			Age:         ageparam,
			VaccineType: vaccinetypeparam,
		},
	}

	fmt.Println("Enviando datos al servidor")
	res, err := c.RegUser(context.Background(), request)

	if err != nil {
		log.Fatal("Error en enviar peticion %v", err)
	}

	fmt.Println("Todo good: ", res.Resultado)

}

func http_server(w http.ResponseWriter, r *http.Request) {
	instance_name := os.Getenv("NAME")
	//instance_name := "grpcInstancia"
	fmt.Println("Manejando peticion HTTP cliente: ", instance_name)

	if r.URL.Path != "/" {
		http.Error(w, "404 No encontrado.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Println("Raiz de HTTP para cliente")
		http.StatusText(http.StatusAccepted)

	case "POST":
		fmt.Println("Iniciando envio de mensajes")
		decoder := json.NewDecoder(r.Body)

		var us userStruct
		err := decoder.Decode(&us)

		if err != nil {
			fmt.Println("Error al decodificar json de locust: %v", err)
		}

		fmt.Fprintf(w, "Mensaje recibido \n")
		fmt.Fprintf(w, "Nombre es: %s\n", us.Name)

		//enviar el mensaje
		registrarUsuario(us.Name, us.Location, us.Gender, us.Age, us.VaccineType)

	default:
		fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
		return
	}
}

func main() {
	//instance_name := os.Getenv("NAME")
	instance_name := "grpcInstancia"
	//client_host := os.Getenv("CLIENT_HOST")
	client_host := ":8080"

	fmt.Println("Cliente ", instance_name, " listo!")
	fmt.Println("Iniciando http server en ", client_host)

	http.HandleFunc("/", http_server)

	if err := http.ListenAndServe(client_host, nil); err != nil {
		log.Fatal(err)
	}

}
