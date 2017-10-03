package DataStream

import (
	"fmt"
	"testing"
)

func TestNewDataStream(t *testing.T) {
	dataStream := NewDataStream("passballBean.bytes")
	fmt.Println(dataStream.data)
	id, _ := dataStream.ReadFloat32()
	name, _ := dataStream.ReadUTF()
	fmt.Println(id)
	fmt.Println(string(name))
}
