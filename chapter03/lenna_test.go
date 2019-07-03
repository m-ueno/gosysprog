package main_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
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
	if bytes.Equal(buffer, []byte("tEXt")) {
		rawText := make([]byte, length)
		chunk.Read(rawText)
		fmt.Println(string(rawText))
	}
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

// 3.5.4
func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer bytes.Buffer // これはポインタではない (p.45)

	// 以下 length, chunkName, data, CRC の順に書き込む
	// 文字列とバイト列はそのまま書きこむけど
	// 数値はbig endianで書きこむ

	binary.Write(&buffer, binary.BigEndian, int32(len(byteData)))
	buffer.WriteString("tEXt")
	buffer.Write(byteData)
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	binary.Write(&buffer, binary.BigEndian, crc.Sum32())

	return &buffer
}

func ExampleTextEmbedding() {
	file, err := os.Open("Lenna.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newFile, err := os.Create("Lenna2.png")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	chunks := readChunks(file)

	io.WriteString(newFile, "\x89PNG\r\n\x1a\n") // PNGファイルシグネチャ (8bytes)
	io.Copy(newFile, chunks[0])
	io.Copy(newFile, textChunk("ASCII PROGRAMMING++"))
	for _, chunk := range chunks[1:] {
		io.Copy(newFile, chunk)
	}

	newFile.Close()

	newFileR, _ := os.Open("Lenna2.png")
	defer newFileR.Close()
	for _, chunk := range readChunks(newFileR) {
		dumpChunk(chunk)
	}

	// Output:
	// chunk 'IHDR' (13 bytes)
	// chunk 'tEXt' (19 bytes)
	// ASCII PROGRAMMING++
	// chunk 'sRGB' (1 bytes)
	// chunk 'IDAT' (473761 bytes)
	// chunk 'IEND' (0 bytes)
}
