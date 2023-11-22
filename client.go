package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	addrc  = "localhost:12345"
	protoc = "tcp4"
)

func main() {
	// Подключение к сетевой службе.
	conn, err := net.Dial(protoc, addrc)
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем подключение.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	reader := bufio.NewReader(conn)
	id := 0

	go func() {
		for {
			pvb, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			id++
			str := strings.Trim(string(pvb), "\n")
			str = strings.Trim(str, "\r")
			fmt.Printf("Найдена поговорка № %d: %s\n", id, str)
		}
	}()

	fmt.Println("Для выхода из программы введите: Выход")
	s := ""
	for {
		_, err := fmt.Scanln(&s)
		if err != nil {
			return
		}
		switch s {
		case "Выход":
			log.Println("Выход из программы")
			return
		}
	}
}
