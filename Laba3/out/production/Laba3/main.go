package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	matchSem   = make(chan int, 1) // Канал для сигналу "сірники на столі"
	tobaccoSem = make(chan int, 1) // Канал для сигналу "тютюн на столі"
	paperSem   = make(chan int, 1) // Канал для сигналу "папір на столі"
	smokerSem  = make(chan int, 0) // Канал для сигналу курцям, що можна скручувати цигарку
	agentSem   = make(chan int, 0) // Канал для сигналу посереднику, що можна поставити інгредієнти на стіл
	waitGroup  sync.WaitGroup      // Об'єкт для синхронізації закінчення всіх потоків
)

func agent() {
	for {
		<-agentSem // Очікуємо на сигнал від посередника
		rand.Seed(time.Now().UnixNano())
		ingredients := rand.Intn(3) // Вибираємо випадковий набір інгредієнтів

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

		smokerSem <- 1 // Сповіщуємо курців, що можна скручувати цигарку
	}
}

func smoker(name string, ingredientSem chan int, otherSem chan int) {
	for {
		<-ingredientSem // Очікуємо на сигнал про наявність інгредієнта
		fmt.Printf("%s бере інгредієнт і скручує цигарку.\n", name)
		time.Sleep(time.Second) // Моделюємо час на скручування цигарки
		fmt.Printf("%s курить цигарку.\n", name)
		time.Sleep(time.Second) // Моделюємо час на куріння цигарки
		fmt.Printf("%s закінчив куріння.\n", name)

		otherSem <- 1 // Сповіщуємо посередника, що інгредієнт використано
	}
}

func main() {
	waitGroup.Add(4) // Ініціалізуємо об'єкт синхронізації для чотирьох потоків

	go agent() // Запускаємо посередника

	// Запускаємо потоки для курців, передаючи їм канали для інгредієнтів
	go smoker("Сірниковий курець", matchSem, tobaccoSem)
	go smoker("Папірковий курець", tobaccoSem, matchSem)
	go smoker("Тютюновий курець", paperSem, matchSem)

	agentSem <- 1 // Запускаємо посередника, відправляючи перший сигнал

	waitGroup.Wait() // Очікуємо завершення всіх потоків
}
