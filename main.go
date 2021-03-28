package main

import (
	"fmt"
	"log"
	"net"
	"sort"
)

func netWorker(ports, results chan int, address string) {
	// Worker receives up to 100 ports and for the ports in that range it attempts to scan each port in the range
	for p := range ports {
		address := fmt.Sprintf("%v:%d", address, p)
		conn, err := net.Dial("tcp", address)
		// If connection fails it sends a 0 to the results channel to be filtered out
		if err != nil {
			results <- 0
			continue
		}
		// If successful it closes connection and sends port to the results channel
		err2 := conn.Close()
		if err2 != nil {
			log.Fatal("Unable to close connection... Quiting...")
		}

		results <- p
	}
}

func main() {
	// Variable and scans to allow custom selecting of website
	var address string
	fmt.Printf("Enter the web address: ")
	_, err := fmt.Scanf("%v\n", &address)
	if err != nil {
		log.Fatal("A valid web address was not entered...")
	}
	fmt.Print("Scanning... Please Wait...\n")

	// Creating 2 channels, 1 to send ports to our worker and the other to store the results the workers gather
	ports := make(chan int, 100)
	results := make(chan int)
	// Variable to store the ports that are reported as open into a slice
	var openports []int

	// For loop waiting for the ports channel to fill up with 100 items then sends them to the worker
	for i := 0; i < cap(ports); i++ {
		go netWorker(ports, results, address)
	}

	// goRoutine to send ports 1 - 1024 to the ports channel
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// The results the worker sends to the results channel is processed here and stored in the open ports variable
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	// Closing the channels as they are no longer needed
	close(ports)
	close(results)
	// Sorted over the ports reported to the results channel
	sort.Ints(openports)
	// Printing each port in the openports variable
	fmt.Print("The following Ports were detected as open:\n")
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
