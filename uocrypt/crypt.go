package uocrypt

import "bytes"

func Decompress(in chan byte) chan byte {
	o := make(chan byte)
	go func() {
		defer close(o)
		var (
			bitNum  int   = 8
			treePos int32 = 0
			value   int   = 0
			mask    int   = 0
		)
		for {
			if bitNum == 8 {
				b, ok := <-in
				if !ok {
					return
				}
				value = int(b)
				bitNum = 0
				mask = 0x80
			}
			if value&mask != 0 {
				treePos = int32(tree[treePos<<1])
			} else {
				treePos = int32(tree[(treePos<<1)+1])
			}
			mask >>= 1
			bitNum++
			if treePos <= 0 { // leaf
				if treePos == -256 {
					bitNum = 8
					treePos = 0
					continue
				}
				o <- byte(-treePos)
				treePos = 0
			}
		}
	}()
	return o
}

func Compress(d []byte) []byte {
	var buffer bytes.Buffer

	var entryOffset int
	bitNum := uint32(0)
	value := uint32(0)
	offset := 0

	for offset < len(d) {
		entryOffset = int(d[offset]) << 1
		offset++
		bitNum += huffman[entryOffset]
		value <<= huffman[entryOffset]
		value |= huffman[entryOffset+1]

		for bitNum >= 8 {
			bitNum -= 8
			buffer.WriteByte(byte(value >> bitNum))
		}
	}

	bitNum += 4
	value <<= 4
	value |= 0xd

	if bitNum&7 != 0 {
		value <<= (8 - (bitNum & 7))
		bitNum += (8 - (bitNum & 7))
	}
	for bitNum >= 8 {
		bitNum -= 8
		buffer.WriteByte(byte(value >> bitNum))
	}

	return buffer.Bytes()
}
