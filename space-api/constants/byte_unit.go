package constants

type MemoryByteSize uint64

const (
	Byte MemoryByteSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)
