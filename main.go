package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var AppVersion = "Test Version"

func main() {
	fmt.Printf("k8s-test %s started \r\n", AppVersion)

	if len(os.Args) < 2 {
		fmt.Println("Usage: k8s-test <mode> (cpu | inet)")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "cpu":
		cpuLoad()
	case "inet":
		checkConnections()
	default:
		fmt.Println("Invalid mode. Use 'cpu' or 'inet'.")
	}
}

func cpuLoad() {
	fmt.Println("Starting CPU intensive task...")

	// Установим конечное время выполнения задачи
	endTime := time.Now().Add(5 * time.Minute)

	// Получаем количество доступных ядер
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Запускаем воркеры в отдельной горутине для каждого ядра
	var wg sync.WaitGroup
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(endTime) {
				hashRandomData()
			}
		}()
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
	fmt.Println("CPU intensive task finished.")
}

// Функция для генерации и хеширования случайных данных
func hashRandomData() {
	// Создаем случайный массив данных
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	// Вычисляем SHA-256 хеш
	hash := sha256.Sum256(data)

	// Выводим результат (опционально, для предотвращения оптимизаций)
	_ = hash[0]
}

func checkConnections() {
	file, err := os.Open("inet.cfg")
	if err != nil {
		log.Fatalf("Failed to open inet.cfg: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := scanner.Text()
		for i := 0; i < 3; i++ {
			if err := connectToAddress(address); err != nil {
				fmt.Printf("Attempt %d: Failed to connect to %s: %v\n", i+1, address, err)
			} else {
				fmt.Printf("Attempt %d: Successfully connected to %s\n", i+1, address)
			}
			time.Sleep(1 * time.Second)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading inet.cfg: %v", err)
	}
}

func connectToAddress(address string) error {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address format: %s", address)
	}

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
