package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	matchSem   = make(chan int, 1)
	tobaccoSem = make(chan int, 1)
	paperSem   = make(chan int, 1)
	smokerSem  = make(chan int, 0)
	agentSem   = make(chan int, 0)
	waitGroup  sync.WaitGroup
)

func agent() {
	for {
		<-agentSem
		rand.Seed(time.Now().UnixNano())
		ingredients := rand.Intn(3)

		switch ingredients {
		case 0:
			fmt.Println("Посередник кладе папір і тютюн на стіл.")
			paperSem <- 1
			tobaccoSem <- 1
		case 1:
			fmt.Println("Посередник кладе сірники і папір на стіл.")
			matchSem <- 1
			paperSem <- 1
		case 2:
			fmt.Println("Посередник кладе сірники і тютюн на стіл.")
			matchSem <- 1
			tobaccoSem <- 1
		}

		smokerSem <- 1
	}
}

func smoker(name string, ingredientSem chan int, otherSem chan int) {
	for {
		<-ingredientSem
		fmt.Printf("%s бере інгредієнт і скручує цигарку.\n", name)
		time.Sleep(time.Second)
		fmt.Printf("%s курить цигарку.\n", name)
		time.Sleep(time.Second)
		fmt.Printf("%s закінчив куріння.\n", name)

		otherSem <- 1
	}
}

func main() {
	waitGroup.Add(4)

	go agent()

	go smoker("Сірниковий курець", matchSem, tobaccoSem)
	go smoker("Папірковий курець", tobaccoSem, matchSem)
	go smoker("Тютюновий курець", paperSem, matchSem)

	agentSem <- 1

	waitGroup.Wait()
}
