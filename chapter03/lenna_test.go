package main_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func dumpChunk(chunk io.Reader) {
	// chunk.Read()とbinary.Read()は違う.

	var length int32
	binary.Read(chunk, binary.BigEndian, &length)
	buffer := make([]byte, 4)
	chunk.Read(buffer)
	fmt.Printf("chunk '%v' (%d bytes)\n", string(buffer), length)
}

func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader

	file.Seek(8, 0)
	var offset int64 = 8
	for {
		var length int32 // 32bit (4bytes) 読む
		err := binary.Read(file, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}

		// なぜ+12?
		// binary.Readではlengthの4bytesはまだ消化されていないのか.
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))

		// チャンク名(4bytes) + データ長 + CRC(4bytes) 先に移動
		offset, _ = file.Seek(int64(length+8), 1)
	}

	return chunks
}

func ExampleMain() {
	file, err := os.Open("Lenna.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	chunks := readChunks(file)
	for _, chunks := range chunks {
		dumpChunk(chunks)
	}

	// Output:
	// chunk 'IHDR' (13 bytes)
	// chunk 'sRGB' (1 bytes)
	// chunk 'IDAT' (473761 bytes)
	// chunk 'IEND' (0 bytes)
}
