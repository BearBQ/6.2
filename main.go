package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func randomNum(min int, max int) int {
	result := min + rand.Intn(max-min+1)
	return result
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-c
		log.Println("Получен сигнал остановки")
		cancel()
	}()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := -1000
	max := 1000
	var wg sync.WaitGroup
	for i := 1; i < 50; i++ {
		wg.Add(1)
		random := randomNum(min, max)
		go worker(ctx, &wg, random, i)
		time.Sleep(50 * time.Millisecond)

	}
	wg.Wait()
	fmt.Println("Program is finished")
}

func worker(ctx context.Context, wg *sync.WaitGroup, num int, id int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Паника:", r, "Число: ", num)
		}
	}()
	defer wg.Done()
	select {
	case <-ctx.Done():
		log.Printf("Worker %d получил сигнал остановки", id)
		return
	default:
		if num <= 0 {
			panic("число меньше или равно нулю")
		}
		fmt.Println(num * num)
	}
}
