package structs
import (
	'assert'
	'binary'
)

type BNode struct {
	data []byte
}

const (
	BNODE_NODE = 1
	BNODE_LEAF= 2
)

type BTree struct {
	root uint64
	// callacks for on-disk pages management
	get func(uint64) BNode
	new func(BNode) uint64 // new page allocation
	del func(uint64)
}

const HEADER = 4
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init(){
	nodemax:= HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(nodemax <= BTREE_PAGE_SIZE)
}

// helper functions for header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}
func(node BNode) setHeader(btype uint16,nkeys uint16){
	binary.LittleEndian.PutUint16(node.data[0:2],btype)
	binary.LittleEndian.PutUint16(node.data[2:4],nkeys)
}

// pointes
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.uint64(node.data[pos:])
}

func (node BNode) setPtr(idx uint16,val uint64) {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:],val)
}

func offsetPos(node BNode, idx uint16) int16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node,idx):])
}
func (node BNode) setOffset(idx uint16,offst uint16){
	binary.LittleEndian.PutUint16(node.data[offsetPos(node,idx):], offset)
}
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}
func (node BNode) getKey(idx uint16) [] byte{ 
	assert(idx<node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}
func (node BNode) getVal(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

func (node BNode) sizeNBytes() uint16{
	return node.kvPos(node.nkeys())
}

func bNodeLookLE(node BNode, key []byte) uint16{
	nkeys := node.nkeys()
	found := uint16(0)
	for i := uint16(1); i< nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i),key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break 
		}
	}
	return found
}