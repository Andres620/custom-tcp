package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"image"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	ownFunctions "swe-challenge-tcp/functions"

	"github.com/pkg/errors"
)

// structs
type imagePNG struct {
	Dx  int
	Dy  int
	Pix []uint8
}

type ImageGOB struct {
	Dx  int
	Dy  int
	Pix []uint8
}

// puerto del server - no se si deba cambiar el puerto para los canales
const (
	Port = ":61000"
)

type HandleFunc func(net.Conn, *bufio.ReadWriter)

type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc
	mapLock  sync.Mutex // Los mapas no son seguros para subprocesos, por lo que se necesitan bloqueos de exclusión mutua para controlar el acceso.
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

func (endPoint *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	endPoint.mapLock.Lock()
	endPoint.handler[name] = f
	endPoint.mapLock.Unlock()
}

// se encarga de manjear los mensajes y llamar el comando apropiado
func (endPoint *Endpoint) handleMessages(conn net.Conn) {
	// Envuelva la conexión al lector de búfer para una fácil lectura
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Leer desde la conexión hasta que se encuentre EOF Espere que la siguiente entrada sea el nombre del comando. Llame al controlador registrado para el comando.
	log.Print("Received commands")
	for {
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n  ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}

		// Recortar los retornos de carro y los espacios adicionales en la cadena de solicitud-ReadString no eliminará ninguna nueva línea.
		cmd = strings.Trim(cmd, "\n ")
		log.Println("cmd  -> " + cmd)
		//log.Println("param-> " + param)
		// Obtenga la función de manejador apropiada del mapeo de manejador y llámela.
		endPoint.mapLock.Lock()
		handleCommand, ok := endPoint.handler[cmd]
		endPoint.mapLock.Unlock()
		if !ok {
			//SI EL COMANDO ENTRA AQUI, SE APAGA EL SERVER Y NO RESPONE, CORREGIR
			log.Println("Command '" + cmd + "' is not registered. -Server")
			return
		} else {
			handleCommand(conn, rw)
		}

	}
}

func handleStrings(conn net.Conn, rw *bufio.ReadWriter) {
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}

	s = strings.Trim(s, "\n ")
	log.Println("STRING IN HANDLESTRINGS: ", s)

	_, err = rw.WriteString("ANS TO STRING COMMAND\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}

	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func handleImageGOB(conn net.Conn, rw *bufio.ReadWriter) {
	var data ImageGOB
	// Envuelva la conexión al lector de búfer para una fácil lectura
	//r := bufio.NewReader(strings.NewReader(param))

	err := gob.NewDecoder(rw).Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	fmt.Println("GOB dx:", data.Dx)
	fmt.Println("GOB dy:", data.Dy)

	rect := image.Rect(0, 0, data.Dx, data.Dy)
	newImg := ownFunctions.BuildImage(rect, data.Pix)
	ownFunctions.Save("server/serverfiles/newImg.png", newImg)

	_, err = rw.WriteString("GOB DATA IMG WAS CONVERTED INTO AN IMAGE\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}

	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func (endPoint *Endpoint) Listen() error {
	var err error
	endPoint.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrap(err, "Unable to listen on "+Port+"\n")
	}
	log.Println("Listen on", endPoint.listener.Addr().String())
	for {
		conn, err := endPoint.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go endPoint.handleMessages(conn)
	}
}

func Server() error {
	endpoint := NewEndpoint()

	// Agregar función de controlador

	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("IMAGEGOB", handleImageGOB)
	// empieza a escuchar

	return endpoint.Listen()
}
