package main

import (
	"github.com/Ivan-Bolotov/golang-calculating-project/internal/servers"
	"sync"
	"time"
)

var wg sync.WaitGroup // Для синхронизации

func main() {
	// Старт основного сервера (оркестратора) на порту из переменной окружения в .env файле
	wg.Add(1)
	go func() {
		servers.StartNewHttpStorageServer()
		wg.Done()
	}()
	// Немного ждём для запуска оркестратора
	time.Sleep(time.Millisecond * 100)

	// Здесь можем запустить сколко угодно серверов вычислителей в разных горутинах (по умолчанию один сервер на 8081 порту)
	for i := 0; i < 1; i++ {
		go servers.StartNewHttpComputingServer(8081 + i) // Передаём порт для запуска
	}

	wg.Wait() // Ожидаем завершения работы оркестратора
}
