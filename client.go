package main

import (
	"flag"
	"fmt"
	"github.com/weisd/goicomet/client"
	"os"
)

const (
	SUBURL  = "http://localhost:8100/poll"
	PUSHURL = "http://localhost:8000/push"
	SIGNURL = "http://localhost:8000/sign"
)

var (
	cname   = flag.String("c", "test", "test chanal")
	action  = flag.String("act", "sub", "sub action")
	content = flag.String("cont", "content", "push content")
)

func main() {

	flag.Parse()

	out := make(chan bool)

	switch *action {
	case "sub":

		for i := 0; i < 1000; i++ {
			go sub(*cname)
		}

		<-out
	case "push":
		push(*cname, *content)
	}

}

func sub(cname string) {
	SubClient := client.Client{Suburl: SUBURL, Signurl: SIGNURL, Pushurl: PUSHURL, Cname: cname}
	// fmt.Println(SubClient)

	signChan := make(chan bool)
	dataChan := make(chan map[string]interface{})

	go func() {
		for {
			data := SubClient.Sub()
			// fmt.Println("respost: ", data)
			// fmt.Println("-------------")
			// hasData := false
			for i := 0; i < len(data); i++ {
				item := data[i]
				t := item["type"].(string)
				// fmt.Println(t)
				switch t {
				case "401":
					signChan <- true
				case "data":
					// hasData = true
					SubClient.Seq = item["seq"].(float64)
					SubClient.Seq++
					dataChan <- item

					// panic(data)
				default:
					fmt.Println("...")
				}
			}
		}
	}()

	for {
		select {
		case <-signChan:
			// fmt.Println("re sign")
			SubClient.Sign()
		case jsonData := <-dataChan:
			fmt.Println("jsonData ", jsonData)
			// log(jsonData["content"].(string))

		}
	}
}

func push(cname, content string) {
	client := client.Client{Suburl: SUBURL, Signurl: SIGNURL, Pushurl: PUSHURL, Cname: cname}
	// fmt.Println(client)

	client.Push(content)
}

func log(con string) {
	fmt.Println("写log")
	fname := "/Users/weisd/gocode/src/github.com/weisd/goicomet/log.txt"
	f, err := os.OpenFile(fname, os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	n, err := f.WriteString(con + "\r\n")
	if err != nil {
		panic(err)
	}

	fmt.Println(n)
}
