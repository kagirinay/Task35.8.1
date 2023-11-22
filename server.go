package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	pogovorka = "https://go-proverbs.github.io/"
	addrs     = "localhost:12345"
	protos    = "tcp4"
)

func main() {
	str, err := getStr(pogovorka)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск сетевой службы по протоколу TCP
	listener, err := net.Listen(protos, addrs)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			return
		}
	}(listener)

	// Создаёт бесконечный цикл обработки подключений.
	// Для предотвращения завершения обслуживания после первого подключения.
	go func() {
		for {
			// Принимаем подключение
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Установленно соединение:", conn.RemoteAddr())

			// Вызывает обработчик подключения
			go handleConn(conn, str)
		}
	}()

	fmt.Println("Для выхода из программы введите: Выход")
	s := ""
	for {
		_, err := fmt.Scanln(&s)
		if err != nil {
			return
		}
	}
}

// Обработчик. Вызывается для каждого созданного соединения.
func handleConn(conn net.Conn, str []string) {
	// Отложенное закрытие соединения.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// Выборочное считывание строк 1 раз в 3 секунды.
	for {
		_, err := conn.Write([]byte(str[rand.Intn(len(str))] + "\n\r"))
		if err != nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
}

// Забирает поговорки с сайта.
func getStr(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := parseUrl(string(body))
	return data, nil
}

// Выбираем поговорки.
func parseUrl(body string) (data []string) {
	tkn := html.NewTokenizer(strings.NewReader(body))
	var values []string
	var isH3, isPogovorka bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return values
		case tt == html.StartTagToken:
			t := tkn.Token()
			if !isH3 {
				isH3 = t.Data == "h3"
			} else {
				isPogovorka = t.Data == "a"
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if isPogovorka {
				values = append(values, t.Data)
			}
			isH3 = false
			isPogovorka = false
		}
	}
}
