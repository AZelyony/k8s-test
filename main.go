package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: app <mode> (cpu | inet)")
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

	endTime := time.Now().Add(5 * time.Minute)
	for time.Now().Before(endTime) {
		findPrime()
	}

	fmt.Println("CPU intensive task finished.")
}

// Ресурсоемкая задача для нагрузки на CPU (поиск простых чисел)
func findPrime() {
	_, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}
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
