package main

import (
	"fmt"
	"math/rand"
	"time"
)

func hello(msg string) {
	fmt.Println(msg + " - goroutine")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Second)
}

func main() {
	go hello("Savio")
	go hello("Cintia")
	time.Sleep(time.Duration(1 * time.Second))
	fmt.Println("Chamada nomal")
}
