package client

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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

func open(addr string) (*bufio.ReadWriter, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}

	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// Llamado cuando la aplicación usa -connect = dirección ip

func Client(ip string) error {

	rw, err := open(ip + Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
	}

	for {
		fmt.Print(">> ")
		//TYPE COMAND
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		cmd, param := ownFunctions.ParseCommand(command)
		fmt.Println("1- cmd: ", cmd)
		fmt.Println("1- param: ", param)

		// Disconect the client
		if strings.TrimSpace(string(command)) == "STOP" {
			log.Println("TCP client exiting...")
			return nil
		}

		// Handle strings
		if strings.TrimSpace(string(cmd)) == "STRING" {
			// Send the command to the server
			n, err := rw.WriteString(cmd + "\n")
			if err != nil {
				return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
			}
			n, err = rw.WriteString(param + "\n")
			if err != nil {
				return errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
			}
		}

		// Handle ImageGOB
		if strings.TrimSpace(string(cmd)) == "IMAGEGOB" {
			dx, dy, pix := ownFunctions.GetImageInfo(strings.TrimSpace(string(param)))
			testStruct := ImageGOB{
				Dx:  dx,
				Dy:  dy,
				Pix: pix,
			}
			fmt.Printf("Data: \n%#v\n", testStruct.Dx)
			fmt.Printf("Data: \n%#v\n", testStruct.Dy)

			// Send the command to the server

			enc := gob.NewEncoder(rw)
			n, err := rw.WriteString(cmd + "\n")
			if err != nil {
				return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
			}
			err = enc.Encode(testStruct)

			if err != nil {
				return errors.Wrap(err, "Encode failed for struct")
			}
		}

		// Flush the buffer
		err = rw.Flush()
		if err != nil {
			return errors.Wrap(err, "Flush failed.")
		}

		// leer respuesta
		response, err := rw.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Command '"+cmd+"' is not registered. -client")
			//log.Println("Command '" + command + "' is not registered. -client")
		} else {
			log.Println("->", response)
		}

	}
}

//IMAGEGOB client/clientfiles/test.jpg
//IMAGEGOB client/clientfiles/img1.jpg
