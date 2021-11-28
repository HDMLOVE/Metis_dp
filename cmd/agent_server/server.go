package main

import "fmt"
import "net"
import "bufio"

// handle the connect function
func process(conn net.Conn) {
	//close the connect
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		// read data
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read data from client err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("receive from client data:", recvStr)
		conn.Write([]byte(recvStr))
	}

	return
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		fmt.Println("listen on port:20000....")
		// waiting for  client to connect
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen failed, err:", err)
			continue
		}
		// start a goroutine to process the connection
		go process(conn)
	}

	return
}
