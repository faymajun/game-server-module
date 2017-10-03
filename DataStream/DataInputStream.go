package DataStream

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"unicode/utf8"
)

type DataStream struct {
	data  []byte
	index int
	len   int
}

func NewDataStream(filePath string) *DataStream {
	data, err := ioutil.ReadFile(filePath)
	check(err)
	return &DataStream{data, 0, len(data)}
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (d *DataStream) ReadFloat32() (float32, error) {
	count, err := d.ReadUInt32()
	check(err)
	return math.Float32frombits(count), nil
}

func (d *DataStream) ReadUInt32() (uint32, error) {
	buff, error := d.readToBuff(4)
	check(error)
	if len(buff) < 4 {
		return 0, errors.New(" length is small.")
	}
	if (buff[0] | buff[1] | buff[2] | buff[3]) < 0 {
		return 0, errors.New(" byte is error.")
	}
	return uint32(uint32(buff[0])<<24 + uint32(buff[1])<<16 + uint32(buff[2])<<8 + uint32(buff[3])<<0), nil
}

func (d *DataStream) ReadInt32() (int32, error) {
	buff, error := d.readToBuff(4)
	check(error)
	if len(buff) < 4 {
		return 0, errors.New(" length is small.")
	}
	if (buff[0] | buff[1] | buff[2] | buff[3]) < 0 {
		return 0, errors.New(" byte is error.")
	}
	return int32(int(buff[0])<<24 + int(buff[1])<<16 + int(buff[2])<<8 + int(buff[3])<<0), nil
}

func (d *DataStream) ReadInt16() (int16, error) {
	buff, error := d.readToBuff(2)
	check(error)
	if len(buff) < 2 {
		return 0, errors.New(" length is small.")
	}
	if (buff[0] | buff[1]) < 0 {
		return 0, errors.New(" byte is error.")
	}
	return int16(int(buff[0])<<8 + int(buff[1])<<0), nil
}

func (d *DataStream) ReadUTF() ([]rune, error) {
	count, error := d.ReadInt16()
	len := int(count)
	check(error)
	buff, error := d.readToBuff(len)
	var runes []rune
	index := 0
	for index < len {
		rune, size := utf8.DecodeRune(buff[index:len])
		if size <= 0 {
			return nil, errors.New(" byte is error.")
		}
		runes = append(runes, rune)
		index += size
	}
	return runes, nil
}

func (d *DataStream) Available() bool {
	if d.len-d.index > 1 {
		return true
	} else {
		return false
	}
}

func (d *DataStream) readToBuff(count int) ([]byte, error) {
	if count > d.len-d.index {
		return nil, errors.New("count over range.")
	}
	ret := d.data[d.index : d.index+count]
	d.index += count
	return ret, nil
}
