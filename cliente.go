package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func main() {
	go Cliente()
	var input string
	fmt.Scanln(&input)
}

func Cliente() {
	con, error := net.Dial("tcp", ":5000")
	if error != nil {
		fmt.Println(error)
		return
	}
	var datos [2]uint64
	error = gob.NewDecoder(con).Decode(&datos)
	if error != nil {
		fmt.Println(error)
	} else {
		canal := make(chan uint64)
		go ProcesoC(datos[0], datos[1], canal)
		for {
			datos[1] = <-canal
			error := gob.NewEncoder(con).Encode(datos[1])
			if error != nil {
				fmt.Println(error)
				return
			}
		}
	}
}

func ProcesoC(x uint64, cont uint64, canal chan uint64) {
	for {
		fmt.Println(x, " = ", cont)
		cont = cont + 1
		canal <- cont
		time.Sleep(time.Millisecond * 500)
	}
}
