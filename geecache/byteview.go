package geecache

type ByteView struct {
	bytes []byte
}

func (b *ByteView) Len() int {
	return len(b.bytes)
}

func (b *ByteView) ByteSlice() []byte {
	return cloneBytes(b.bytes)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (b *ByteView) String() string {
	return string(b.bytes)
}
