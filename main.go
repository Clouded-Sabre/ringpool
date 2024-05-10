package main

import (
	"fmt"
	"time"

	"github.com/Clouded-Sabre/ringpool/lib"
)

// ConcreteData represents concrete data type
type ConcreteData struct {
	content []byte
}

func (c *ConcreteData) SetContent(s string) {
	c.content = []byte(s)
}

// Reset resets the content of the concrete data
func (c *ConcreteData) Reset() {
	fmt.Print("")
}

// PrintContent prints the content of the concrete data
func (c *ConcreteData) PrintContent(length int) {
	fmt.Println("Content:", string(c.content[:length]))
}

func main() {
	lib.Debug = true
	// Define a function that creates a new instance of ConcreteData
	newData := func(length int) lib.DataInterface {
		return &ConcreteData{
			content: make([]byte, length),
		}
	}

	// Create a new RingPool instance
	pool := lib.NewRingPool(10, 100, newData)

	// Get an element from the pool
	element := pool.GetElement(newData)
	fmt.Println("The number of available element is", pool.AvailableChunks())

	// Use the element
	s := "Hohoho"
	element.Data.(*ConcreteData).SetContent(s)
	element.Length = len(s)
	element.Data.PrintContent(element.Length)

	// Return the element to the pool
	pool.ReturnElement(element)
	fmt.Println("The number of available element is", pool.AvailableChunks())

	time.Sleep(20 * time.Second)
}
