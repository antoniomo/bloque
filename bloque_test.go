package bloque

import (
	"fmt"
	"testing"
)

const (
	pushAmount = 100
)

func TestPushBackItems(t *testing.T) {
	// Default block size is 64
	b := New()

	for i := 0; i < pushAmount; i++ {
		newPage := b.PushBackItem(i)
		if newPage == true && (i+1)%defaultBlockSize != 0 {
			t.Errorf("pushing created a new page, but we were inserting item %d, blockSize %d",
				i, defaultBlockSize)
			break
		}
		if newPage == false && (i+1)%defaultBlockSize == 0 {
			t.Errorf("pushing didn't create a new page, but we were inserting item %d, blockSize %d",
				i, defaultBlockSize)
			break

		}
		if have, expect := b.Len(), i+1; have != expect {
			t.Errorf("BloQue.Len() returned %d, expecting %d", have, expect)
			break
		}
		if have, expect := b.NumBlocks(), (i/defaultBlockSize)+1; have != expect {
			t.Errorf("BloQue.NumBlocks() returned %d, expecting %d", have, expect)
			break
		}
	}
}

func TestPopBackItems(t *testing.T) {
	b := New()

	if _, ok := b.PopBackItem(); ok {
		t.Errorf("ok when popping from empty BloQue, BloQue.Len() %d", b.Len())
	}

	// Not a real unit test if we depend on this other method but if it
	// passed the PushBackItem test, we are good :D
	for i := 0; i < pushAmount; i++ {
		b.PushBackItem(i)
	}

	for i := pushAmount - 1; i >= 0; i-- {
		item, ok := b.PopBackItem()
		if !ok {
			t.Errorf("not ok when popping element %d, BloQue.Len() %d", i, b.Len())
			break
		}
		intItem, ok := item.(int)
		if !ok {
			t.Error("error on type assertion from popped item, item:", item)
			break
		}
		if have, expect := intItem, i; have != expect {
			t.Errorf("popped item doesn't match, have %d, expecting %d", have, expect)
			break
		}
	}
	if have, expect := b.Len(), 0; have != expect {
		t.Errorf("wrong length on emptied BloQue, BloQue.Len() %d, expecting %d", have, expect)
	}
	if have, expect := b.NumBlocks(), 1; have != expect {
		t.Errorf("wrong number of blocks on emptied BloQue, BloQue.NumBlocks() %d, expecting %d", have, expect)
	}
	if _, ok := b.PopBackItem(); ok {
		t.Errorf("ok when popping from empty BloQue, BloQue.Len() %d", b.Len())
	}
}

func TestPeekBackItems(t *testing.T) {
	b := New()

	if _, ok := b.PeekBackItem(); ok {
		t.Errorf("ok when peeking from empty BloQue, BloQue.Len() %d", b.Len())
	}

	// Not a real unit test if we depend on PushBackItem but if it passed
	// its test, we are good :D
	for i := 0; i < pushAmount; i++ {
		b.PushBackItem(i)

		item, ok := b.PeekBackItem()
		if !ok {
			t.Errorf("not ok when peeking element %d, BloQue.Len() %d", i, b.Len())
			break
		}
		intItem, ok := item.(int)
		if !ok {
			fmt.Println(intItem)
			t.Errorf("error on type assertion from peeked item")
		}
		if have, expect := intItem, i; have != expect {
			t.Errorf("peeked item doesn't match, have %d, expecting %d", have, expect)
			break
		}
	}
}

func TestPopFrontBlock(t *testing.T) {
	b := New()

	// Not a real unit test if we depend on this other method but if it
	// passed the PushBackItem test, we are good :D
	for i := 0; i < pushAmount; i++ {
		b.PushBackItem(i)
	}

	// Assuming the final block isn't full
	totalBlocks := (pushAmount / defaultBlockSize)
	for i := 0; i <= totalBlocks; i++ {
		block, ok := b.PopFrontBlock()
		if !ok {
			t.Errorf("not ok when popping front block, BloQue.Len() %d, BloQue.NumBlocks() %d",
				b.Len(), b.NumBlocks())
		}
		expectedItems := defaultBlockSize
		if i == totalBlocks {
			// Last block has the remainder items
			expectedItems = pushAmount % defaultBlockSize
		}
		if have, expect := len(block), expectedItems; have != expect {
			t.Errorf("wrong length on popped block, len(block) %d, expecting %d", have, expect)
		}
		for j, item := range block {
			intItem, ok := item.(int)
			if !ok {
				t.Errorf("error on type assertion from item")
				break
			}
			if have, expect := intItem, j+defaultBlockSize*i; have != expect {
				t.Errorf("item doesn't match, have %d, expecting %d", have, expect)
				break
			}
		}

		if have, expect := b.Len(), pushAmount-(defaultBlockSize*(i+1)); have != expect {
			// If it's not the last block, there should be some
			// items here
			if !(have == 0 && i == totalBlocks) {
				t.Errorf("wrong length on BloQue, BloQue.Len() %d, expecting %d", have, expect)
			}
		}
		if have, expect := b.NumBlocks(), totalBlocks-i; have != expect {
			// If it's not the last block, there should be some
			// more blocks here
			if !(have == 1) && i == totalBlocks {
				t.Errorf("wrong number of blocks on BloQue, BloQue.NumBlocks() %d, expecting %d", have, expect)
			}
		}
	}
	if have, expect := b.Len(), 0; have != expect {
		t.Errorf("wrong length on BloQue, BloQue.Len() %d, expecting %d", have, expect)
	}
	if have, expect := b.NumBlocks(), 1; have != expect {
		t.Errorf("wrong number of blocks on BloQue, BloQue.NumBlocks() %d, expecting %d", have, expect)
	}
}
