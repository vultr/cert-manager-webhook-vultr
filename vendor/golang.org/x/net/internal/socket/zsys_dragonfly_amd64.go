// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs defs_dragonfly.go

package socket

type iovec struct {
	Base *byte
	Len  uint64
}

type msghdr struct {
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *iovec
	Iovlen     int32
	Pad_cgo_1  [4]byte
	Control    *byte
	Controllen uint32
	Flags      int32
}

type cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

const (
	sizeofIovec  = 0x10
	sizeofMsghdr = 0x30
)
