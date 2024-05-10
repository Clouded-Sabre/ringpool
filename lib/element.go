package lib

import (
	"fmt"
	"time"
)

// Payload represents a single chunk of payload
type DataInterface interface {
	Reset()
	PrintContent(length int)
}

type Element struct {
	Data           DataInterface
	Length         int // length of the real data
	index          int
	LastAllocation time.Time
	CallStack      []string
	StayAtChannel  string
}

type NewData func(length int) DataInterface

// NewChunk creates a new chunk with the given length
func NewElement(index, length int, newData NewData) *Element {
	return &Element{
		Data:   newData(length),
		Length: 0,
		index:  index,
		//LastAllocation: time.Time{},
		CallStack:     nil,
		StayAtChannel: "",
	}
}

// AddCallStack adds a function string to the call stack of the chunk
func (e *Element) AddCallStack(funcString string) {
	e.CallStack = append(e.CallStack, funcString)
}

// PopCallStack removes the last function string from the call stack of the chunk
func (e *Element) PopCallStack() {
	if len(e.CallStack) > 0 {
		e.CallStack = e.CallStack[:len(e.CallStack)-1]
	}
}

// AddToChannel assign a channel string to StayAtChannel of the chunk
func (e *Element) AddToChannel(channelString string) {
	e.StayAtChannel = channelString
}

func (e *Element) RemoveFromChannel() {
	e.StayAtChannel = ""
}

// PrintCallStack prints the call stack of the chunk
func (e *Element) PrintCallStack() {
	fmt.Print("Call Stack:")
	for _, call := range e.CallStack {
		fmt.Printf(" -> %s", call)
	}
	fmt.Println()
	if len(e.StayAtChannel) > 0 {
		fmt.Println("Chunk@channel:", e.StayAtChannel)
	}
	fmt.Println()
	e.Data.PrintContent(e.Length)
}

// Reset resets necessary fields after chunk is returned
func (e *Element) Reset() {
	e.Length = 0
	e.LastAllocation = time.Time{}
	e.CallStack = nil
	e.Data.Reset()
}

// IsTimedOut checks if an allocated chunk is been held for too long
func (e *Element) IsTimedOut(timeout time.Duration) bool {
	return time.Since(e.LastAllocation) > timeout
}

// GetAgeSeconds returns how long has it been since the chunk was last allocated
func (e *Element) GetAgeSeconds() float64 {
	return time.Since(e.LastAllocation).Seconds()
}
