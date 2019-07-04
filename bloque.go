package bloque

import (
	"container/list"
)

const (
	// A sane default, see
	// https://github.com/juju/utils/blob/b297061d1e5ae1b0d71c02799238b05299205fa2/deque/deque.go#L33-L37
	defaultBlockSize = 64
)

// Option ...
type Option func(*BloQue)

// BlockSize ...
func BlockSize(blockSize int) Option {
	return func(blk *BloQue) {
		blk.blockSize = blockSize
	}
}

// BloQue represents a block queue with push/pop items on the back interface
//
// It can be thought of as an items stack (as you push items on the back) while
// allowing FIFO semantics for blocks (that is, it is a block queue allowing you
// to pop out the first block). See the README for use cases.
//
// Code is inspired in the deque implementation at juju:
// https://github.com/juju/utils/blob/master/deque/deque.go
type BloQue struct {
	blockSize  int
	maxLength  int
	maxBlocks  int
	usedLength int
	backIdx    int
	blocks     list.List
}

// BlockT is the block type
type BlockT []interface{}

func newDefault() *BloQue {
	// Default options
	blk := &BloQue{
		blockSize: defaultBlockSize,
		maxLength: 0, // No limit
		maxBlocks: 0, // No limit
		backIdx:   0,
	}
	return blk
}

// New returns a new BloQue with any options
func New(opt ...Option) *BloQue {

	// Default options
	blk := newDefault()

	for _, setter := range opt {
		setter(blk)
	}

	blk.blocks.PushBack(blk.newBlock())

	return blk
}

// NewWithMaxLength returns a new BloQue with a length (items) cap
func NewWithMaxLength(maxLength int, opt ...Option) *BloQue {

	// Default options
	blk := newDefault()

	for _, setter := range opt {
		setter(blk)
	}

	blk.maxBlocks = blk.maxLength / blk.blockSize
	blk.blocks.PushBack(blk.newBlock())

	return blk

}

// NewWithMaxBlocks returns a new BloQue with a length (blocks) cap
func NewWithMaxBlocks(maxBlocks int, opt ...Option) *BloQue {

	// Default options
	blk := newDefault()

	for _, setter := range opt {
		setter(blk)
	}

	blk.maxLength = blk.maxBlocks * blk.blockSize
	blk.blocks.PushBack(blk.newBlock())

	return blk
}

func (b *BloQue) newBlock() BlockT {
	return make(BlockT, b.blockSize)
}

// Len returns the amount of items on a BloQue
func (b *BloQue) Len() int {
	return b.usedLength
}

// NumBlocks returns the amount of blocks on a BloQue
//
// Note that an empty BloQue has 1 block allocated and ready to use
func (b *BloQue) NumBlocks() int {
	return b.blocks.Len()
}

// PushBackItem ...
func (b *BloQue) PushBackItem(item interface{}) bool {

	// If the current block is full, allocate the next one on write
	if b.backIdx == b.blockSize {
		b.blocks.PushBack(b.newBlock())
		b.backIdx = 0
		// TODO: maxBlocks cap
	}

	block := b.blocks.Back().Value.(BlockT)
	block[b.backIdx] = item
	b.backIdx++
	b.usedLength++
	// TODO: maxLength cap

	// If we completed a block, return true
	return (b.backIdx == b.blockSize)
}

// PopBackItem ...
func (b *BloQue) PopBackItem() (interface{}, bool) {

	if b.usedLength < 1 {
		return nil, false
	}

	listBlk := b.blocks.Back()
	blk := listBlk.Value.(BlockT)
	b.backIdx--
	item := blk[b.backIdx]
	blk[b.backIdx] = nil
	b.usedLength--

	// Remove empty block unless it's the only one
	if b.backIdx == 0 && b.blocks.Len() > 1 {
		b.blocks.Remove(listBlk)
		b.backIdx = b.blockSize
	}

	return item, true
}

// PeekBackItem ...
func (b *BloQue) PeekBackItem() (interface{}, bool) {

	if b.usedLength < 1 {
		return nil, false
	}

	item := b.blocks.Back().Value.(BlockT)[b.backIdx-1]
	return item, true
}

// PopFrontBlock ...
func (b *BloQue) PopFrontBlock() (BlockT, bool) {

	if b.usedLength < 1 {
		return nil, false
	}

	listBlk := b.blocks.Front()
	blk := listBlk.Value.(BlockT)
	if b.blocks.Len() == 1 {
		blk = blk[:b.backIdx] // Ensure the returned slice has the right length
	}
	b.blocks.Remove(listBlk)

	if b.blocks.Len() == 0 {
		// We were only using one block, BloQue is now empty
		b.backIdx = 0
		b.usedLength = 0
		b.blocks.PushBack(b.newBlock())
	} else {
		// After removing the front block, recalculate length
		b.usedLength -= b.blockSize
	}

	return blk, true
}
