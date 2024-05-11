package lib

import (
	"log"
	"sync"
	"time"
)

var Debug bool

// GenericPayloadPool is a generic implementation of PayloadPool
type RingPool struct {
	chunks          []*Element
	capacity        int
	readIdx         int
	writeIdx        int
	isFull, isEmpty bool
	allocatedMap    map[int]*Element
	mtx             sync.Mutex
	newData         NewData
	DataParams      []interface{}
}

// NewPayloadPool creates a new payload pool with the specified capacity and chunk length
func NewRingPool(capacity int, newData NewData, params ...interface{}) *RingPool {
	chunks := make([]*Element, capacity)
	for i := 0; i < capacity; i++ {
		chunks[i] = NewElement(i, newData, params...)
	}

	p := &RingPool{
		chunks:       chunks,
		capacity:     capacity,
		allocatedMap: make(map[int]*Element),
		newData:      newData,
		DataParams:   params,
	}

	// start timeout checks for allocated chunks
	if Debug {
		go p.CheckTimedOutChunks()
	}

	return p
}

// GetPayload retrieves a payload from the pool
func (p *RingPool) GetElement() *Element {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Check if the pool is empty
	if p.isEmpty {
		log.Println("Chunk allocation: payload pool is empty, allocate more chunk will impact performance till some chunks are returned")
		return NewElement(p.capacity+1, p.newData, p.DataParams...) // Pool is empty. Create new chunk manually
	}

	chunk := p.chunks[p.readIdx]
	chunk.LastAllocation = time.Now()
	p.readIdx = (p.readIdx + 1) % p.capacity // Move read index circularly

	if p.readIdx == p.writeIdx {
		p.isEmpty = true
	}

	p.isFull = false

	// Add the chunk to allocatedMap
	p.allocatedMap[chunk.index] = chunk

	return chunk
}

// ReturnPayload returns a payload to the pool
func (p *RingPool) ReturnElement(element *Element) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if element.index > p.capacity {
		// manually created chunk, just ignore it
		log.Println("Payload Pool: returned a manually created chunk")
		return
	}

	if p.isFull {
		log.Println("Deallocation: Pool is full, cannot return more chunk")
		return // Pool is full, discard the chunk
	}

	element.Reset()
	p.chunks[p.writeIdx] = element // Reuse the frame object
	// Check if the pool is full
	p.writeIdx = (p.writeIdx + 1) % p.capacity
	if p.writeIdx == p.readIdx {
		p.isFull = true
	}
	p.isEmpty = false

	// remove it from allocatedMap
	delete(p.allocatedMap, element.index)
	element = nil // Set the pointer to nil
}

// AvailableChunks returns the number of available payloads in the pool
func (p *RingPool) AvailableChunks() int {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.readIdx > p.writeIdx {
		return p.capacity - (p.readIdx - p.writeIdx)
	} else if p.readIdx < p.writeIdx {
		return p.writeIdx - p.readIdx
	} else { // ==
		if p.isEmpty {
			return 0
		}
		if p.isFull {
			return p.capacity
		}
		return 0 // put it here just to please compiler hahaha
	}
}

// CheckTimedOutChunks checks every 5 seconds for chunks allocated more than 10 seconds ago
func (p *RingPool) CheckTimedOutChunks() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var count int
	for {
		<-ticker.C
		count = 0

		p.mtx.Lock()
		for _, chunk := range p.allocatedMap {
			if time.Since(chunk.LastAllocation) > 10*time.Second {
				chunk.PrintCallStack()
				count++
			}
		}
		log.Printf("Number of chunks allocated more than 10 seconds ago: %d\n", count)
		p.mtx.Unlock()
	}
}

// Example usage:
/*
func main() {
	pool := NewPayloadPool(2000, 1500) // Adjust capacity and chunk length as needed

	// Example of getting and returning payloads
	payload := pool.GetPayload()
	// Use the payload
	// ...

	// After using the payload, return it to the pool
	pool.ReturnPayload(payload)

	// Check the number of available payloads
	availableChunks := pool.AvailableChunks()
	println("Available chunks:", availableChunks)

	// Close the pool when it's no longer needed
	pool.Close()
}
*/
