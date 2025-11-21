package byte_view

type ByteView struct {
	data []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{
		data: cloneBytes(b),
	}
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.data)
}

func (b ByteView) String() string {
	return string(b.data)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
