package main

import (
	"fmt"
	"log"
	"net"
	//"strconv"
	"strings"
)

type Connections struct {
	connection net.Conn
	file_List  []string
	port       int
}

var counter_slave int = 6000
var counter_client int = 5000

var Slave_List []Connections

func main() {

	fmt.Println("The server is listening on Port 3000")

	listener_Slave, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	listener_Client, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		go acceptLoop(listener_Slave, false)

		acceptLoop(listener_Client, true) // run in the main goroutine
	}
}

func acceptLoop(l net.Listener, check bool) {
	defer l.Close()
	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New connection found!")

		if check == false {
			go listenConnection_Slave(con)
		} else {
			go listenConnection_Client(con)
		}
	}
}

func listenConnection_Slave(conn net.Conn) {

	counter_slave += 1
	output := fmt.Sprintf("%s%d", "Welcome, Slave On port --", counter_slave)

	_, err := conn.Write([]byte(output))
	if err != nil {
		log.Fatalln(err)
	}

	//===================== Split Buffer String ======================

	bufferReader := make([]byte, 100)
	ln1, err := conn.Read(bufferReader)
	if err != nil {
		fmt.Println(err)
	}

	S := string(bufferReader[0:ln1])
	read_Array := strings.Split(S, "--")
	//size := len(read_Array)

	//======================= List of Slaves ==================

	log.Println("Now, Server Schedule the Jobs to Free Slaves")

	new_ip := fmt.Sprintf("%s%d", "localhost:", counter_slave)

	new_conn, err := net.Dial("tcp", new_ip)
	if err != nil {
		// handle error
	} else {
		var obj Connections
		obj.file_List = read_Array
		obj.port = counter_slave
		obj.connection = new_conn

		Slave_List = append(Slave_List, obj)
		log.Println("Slave Port", obj.port, obj.file_List)

		bufferReader := make([]byte, 100)
		ln1, err := new_conn.Read(bufferReader)
		if err != nil {
			fmt.Println(err)
		}

		S := string(bufferReader[0:ln1])

		log.Print(S)
	}
}

func listenConnection_Client(conn net.Conn) {

	counter_client += 1
	output := fmt.Sprintf("%s%d", "Welcome, Client On port --", counter_client)

	_, err := conn.Write([]byte(output))
	if err != nil {
		log.Fatalln(err)
	}

	bufferReader := make([]byte, 100)
	ln1, err := conn.Read(bufferReader)
	if err != nil {
		fmt.Println(err)
	}

	password := string(bufferReader[0:ln1])

	for index, num := range Slave_List {
		output := fmt.Sprintf("%s%d", password, index)
		num.connection.Write([]byte(output))
	}
}
