package lib

import (
	"fmt"
	"log"
	"time"
)

// Payload represents a single chunk of payload
type DataInterface interface {
	Reset()
	PrintContent()
}

type Footprint struct {
	Function  string
	Timestamp time.Time
}

type Element struct {
	Data           DataInterface
	index          int
	LastAllocation time.Time
	footprints     []Footprint
	isAllocated    bool
	pool           *RingPool
}

type NewData func(params ...interface{}) DataInterface

// NewChunk creates a new chunk with the given length
func NewElement(rp *RingPool, index int, newData NewData, params ...interface{}) *Element {
	return &Element{
		Data:  newData(params...),
		index: index,
		pool:  rp,
	}
}

// AddFootPrint adds a function string to the call stack of the element
func (e *Element) AddFootPrint(funcString string) int {
	now := time.Now()
	if len(e.footprints) > 0 {
		lastFootprint := e.footprints[len(e.footprints)-1]
		duration := now.Sub(lastFootprint.Timestamp)
		//threshold := 100 * time.Millisecond // example threshold
		if duration > e.pool.ProcessTimeThreshold {
			log.Printf("Time since last footprint exceeded: %s -> %s, duration: %v\n",
				lastFootprint.Function, funcString, duration)
		}
	}
	e.footprints = append(e.footprints, Footprint{Function: funcString, Timestamp: now})
	return len(e.footprints) - 1
}

// PopCallStack removes the last function string from the call stack of the chunk
func (e *Element) TickFootPrint(pos int) {
	if e.isAllocated {
		if pos >= 0 && pos < len(e.footprints) {
			e.footprints[pos].Function = e.footprints[pos].Function + "✓"
		}
	} else {
		log.Println("Element cannot be ticked because it is already been returned to pool.")
	}
}

// AddChannel assign a channel string to footprints of the chunk
func (e *Element) AddChannel(channelString string) int {
	return e.AddFootPrint("(" + channelString + ")")
}

// AddChannel assign a channel string to footprints of the chunk
func (e *Element) TickChannel() error {
	if e.isAllocated {
		if len(e.footprints) > 0 {
			pos := len(e.footprints) - 1 // the last one
			e.footprints[pos].Function = e.footprints[pos].Function + "✓"
			return nil
		} else {
			return fmt.Errorf("tickChannel: there is no channel to be ticked")
		}
	} else {
		return fmt.Errorf("element cannot be ticked because it is already been returned to pool")
	}
}

// PrintCallStack prints the call stack of the chunk
func (e *Element) PrintCallStack() {
	fmt.Print("Footprints:")
	for _, footprint := range e.footprints {
		fmt.Printf(" -> %s", footprint.Function)
	}
	fmt.Println()
	e.Data.PrintContent()
}

// Reset resets necessary fields after chunk is returned
func (e *Element) Reset() {
	e.LastAllocation = time.Time{}
	e.footprints = nil
	e.isAllocated = false
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
