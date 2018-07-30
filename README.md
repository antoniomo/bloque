# BloQue, a Block Queue of stacked items

Wait, what? What the title means is that this data structure is a block _queue_,
that is, blocks are FIFO (First In First Out), of stacked items, that is, items
are a _stack_ LIFO (Last In First Out). A block is just a bunch of items.

## Why this is useful?

Well there are plenty of other more general data structures. For example a
_deque_ works as both a _queue_ and a _stack_. There are already nice _deque_
implementations in golang, for example this one from the _Juju_ project:
https://github.com/juju/utils/tree/master/deque.

So why do we want this `BloQue` data structure? The use case I found is that I
want to be able to add items one by one to a data structure (list, queue, stack,
whatever) to act like a cache, and once there's enough of them, or after a
pre-specified time has passed, I want to get the first X elements from the front
of that data structure to batch-process them (send to a data pipeline, or
bulk-write to a DB). A normal queue would do, but then I would need to loop over
the first X elements, popping them one at a time. Add in concurrency, and you
need to hold a write lock on that time.

With this `BloQue` data structure you can just add elements one at a time, and
once you are ready to consume them you just pop the front block of the queue.
That's a very fast operation (just touching a few pointers) so it can be done
holding the lock for very little time, and then you can keep adding elements
while processing this front block elsewhere.

The other "extra" stack operations such as Pop and Peek items at the end of the
stack are added for convenience, as well.

I haven't added full _deque_ API on blocks and items because it gets trickier
and I don't see the use case for it, but please raise an issue if needed :)

## Example usage

```go
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/antoniomo/bloque"
)

const (
	elementsPerBlock = 1024
	// Just to show the last block with less than full size
	itemsToProcess = (elementsPerBlock * 10) - 3
)

func consumer(wg *sync.WaitGroup, readyBlocks <-chan bloque.BlockT) {
	// Lets keep a counter of processed blocks
	processedBlocks := 0
	totalProcessedItems := 0

	for blk := range readyBlocks {
		fmt.Printf("Processing BloQue %d, items %d\n", processedBlocks, len(blk))
		processedBlocks++
		totalProcessedItems += len(blk)
		for _, item := range blk {
			// Is this what we expect?
			_, ok := item.(int)
			if !ok {
				log.Panic("oh no, wrong type assertion!")
			}
			// Simulate consumption/processing of the item
			time.Sleep(time.Nanosecond)
		}
	}

	wg.Done()
}

func producer(wg *sync.WaitGroup, readyBlocks chan bloque.BlockT) {

	// Lets prepare a BloQue with a custom elementsPerBlock
	b := bloque.New(bloque.BlockSize(elementsPerBlock))

	for i := 0; i < itemsToProcess; i++ {
		if isCompleted := b.PushBackItem(i); isCompleted {
			// Block completed, send to the processing channel
			frontBlock, _ := b.PopFrontBlock()
			readyBlocks <- frontBlock
		}
	}
	// Last block wasn't completed but lets say we want to process it right
	// away without waiting for more data
	frontBlock, ok := b.PopFrontBlock()
	if !ok {
		log.Panic("ouch, there should be items here")
	}
	readyBlocks <- frontBlock

	close(readyBlocks) // All done, close channel

	wg.Done()
}

func main() {
	fmt.Printf("Processing %d items, block size %d\n", itemsToProcess, elementsPerBlock)

	// Channel to write ready blocks
	readyBlocks := make(chan bloque.BlockT)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go producer(wg, readyBlocks)
	go consumer(wg, readyBlocks)

	wg.Wait()
}

```
