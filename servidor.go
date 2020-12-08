package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func main() {
	var num_canales int = 5
	bloqueados := make([]chan bool, num_canales)
	canal := make([]chan uint64, num_canales)
	var cont uint64 = 0
	for i := 0; i < num_canales; i++ {
		canal[i] = make(chan uint64, 2)
		bloqueados[i] = make(chan bool)
		go Proceso(cont, 0, canal[i], bloqueados[i])
		cont++
	}
	go server(canal, bloqueados)
	var input string
	fmt.Scanln(&input)
}

func server(canal []chan uint64, bloqueados []chan bool) {
	var tam_proceso int = 2
	nuevo_Proceso := make(chan uint64, tam_proceso)
	enviar := make(chan net.Conn)
	resp := make(chan bool)
	con, error := net.Listen("tcp", ":5000")
	if error != nil {
		fmt.Println(error)
		return
	}
	go administrar(canal, bloqueados, nuevo_Proceso, enviar, resp)
	for {
		c, error := con.Accept()

		if error != nil {
			fmt.Println(error)
			continue
		}
		enviar <- c
	}
}

func handleClient(c net.Conn, canal chan uint64, canal_bloq chan bool, nuevo_Proceso chan uint64, resp chan bool) {
	canal_bloq <- true
	datos := [2]uint64{<-canal, <-canal}
	error := gob.NewEncoder(c).Encode(datos)
	if error != nil {
		fmt.Println(error)
	}
	for {
		error := gob.NewDecoder(c).Decode(&datos[1])
		if error != nil {
			nuevo_Proceso <- datos[0]
			nuevo_Proceso <- datos[1]
			resp <- true
			return
		}
	}
}

func Proceso(id uint64, i uint64, canal chan uint64, canal_bloq chan bool) {
	for {
		select {
		case <-canal_bloq:
			canal <- id
			canal <- i
			return
		default:
			fmt.Println(id, " = ", i)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func administrar(canal []chan uint64, bloqueados []chan bool, nuevo_Proceso chan uint64, enviar chan net.Conn, resp chan bool) {
	var tam_proceso int = 2
	for {
		select {
		case <-resp:
			datos := [2]uint64{<-nuevo_Proceso, <-nuevo_Proceso}
			canal = append(canal, make(chan uint64, tam_proceso))
			bloqueados = append(bloqueados, make(chan bool))
			go Proceso(datos[0], datos[1], canal[len(canal)-1], bloqueados[len(bloqueados)-1])
		case c := <-enviar:
			go handleClient(c, canal[0], bloqueados[0], nuevo_Proceso, resp)
			canal = canal[1:]
			bloqueados = bloqueados[1:]
		}
	}
}
