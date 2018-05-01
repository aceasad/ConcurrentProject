package main

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type file_task_avail struct {
	file_List string
	avail     bool
}

var Task_List [5]file_task_avail

func main() {

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		// handle error
	}
	fmt.Println("Connection successful!!", conn.RemoteAddr())

	//===================== Split Buffer String ======================

	bufferReader := make([]byte, 100)
	ln1, err := conn.Read(bufferReader)
	if err != nil {
		fmt.Println(err)
	}

	S := string(bufferReader[0:ln1])
	read_Array := strings.Split(S, "--")
	log.Print(read_Array)
	size := len(read_Array)
	New_Port := read_Array[size-1]

	port, err := strconv.Atoi(New_Port)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	//================= Read the Slave Directory ===============

	folder_name := fmt.Sprintf("%s%d", "../Slave_", port%6)

	file_List := make([]string, 5)

	list := ""

	files, err := ioutil.ReadDir(folder_name)
	if err != nil {
		log.Fatal(err)
	}

	//var i int = 0
	for i, file := range files {
		file_List[i] = file.Name()
		list += file.Name()
		list += "--"
		//i += 1
	}

	//============== Sending To Server to Keep Track ============

	_, err = conn.Write([]byte(list))
	if err != nil {
		log.Fatalln(err)
	}
	//==================== Now Slave Listening ===================
	fmt.Println("I'm Slave, -- listening at --- ", port, "-- Assign Me Task")
	ln, err := net.Listen("tcp", ":"+New_Port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go Searching(conn, file_List, folder_name)
	}

}

func threadworker(threadid int, password string, filename string, done chan string, folder string) {

	fmt.Println("working...")
	//Perform the Search here on filename

	slave_one_file := fmt.Sprintf("%s%s%s", folder, "/", filename)

	_, err := ioutil.ReadFile(slave_one_file)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(slave_one_file)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 5000)

	ln, err := file.Read(buff)
	if err != nil {
		log.Fatal(err)
	}

	file_Data := string(buff[0:ln])

	var status string = "not"
	read_Array := strings.Split(file_Data, "\n")
	for _, element := range read_Array {
		fmt.Println(element)
		if password == element {
			status = "found"
			break
		}
	}

	fmt.Println(" ------------ Searching done--------- ", folder, " : ", filename)

	// Send a value to notify that we're done.

	output := fmt.Sprintf("%d%s", threadid, status)

	done <- output
}

func Searching(conn net.Conn, files []string, folder string) {

	log.Println("Searching")
	output := fmt.Sprintf("%s%d", "Slave make a new connection -- ", conn)
	conn.Write([]byte(output))

	bufferReader := make([]byte, 100)
	ln1, err := conn.Read(bufferReader)
	if err != nil {
		fmt.Println(err)
	}
	SearchPassword := string(bufferReader[0:ln1])

	fmt.Print(SearchPassword)

	Threadnumber := 0
	done := make(chan string)

	fmt.Println(files)

	for _, file := range files {
		go threadworker(Threadnumber, SearchPassword, file, done, folder)
		Threadnumber = Threadnumber + 1
	}
	processdata := []string{<-done, <-done, <-done, <-done, <-done}
	fmt.Println(processdata)
}
