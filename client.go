package main

import (
	//"bufio"
	"fmt"
	//"log"
	"net"
	//"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		// handle error
	} else {
		fmt.Println("Connection successful!!", conn.RemoteAddr())
		bufferReader := make([]byte, 100)
		ln, err := conn.Read(bufferReader)
		if err != nil {
			fmt.Println(err)
		}

		S := string(bufferReader[0:ln])
		fmt.Print(S)
		conn.Write([]byte("berufsverkehr"))
	}

}
