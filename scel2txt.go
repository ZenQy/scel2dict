package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf16"
)

func readUtf16Str(b *bytes.Reader, offset int64, length int) string {
	if offset >= 0 {
		b.Seek(offset, 0)
	}
	data := make([]byte, length)
	b.Read(data)
	u16 := make([]uint16, length/2)
	for i := range u16 {
		u16[i] = binary.LittleEndian.Uint16(data[i*2 : (i+1)*2])
	}
	return string(utf16.Decode(u16))
}

func readUint16(b *bytes.Reader) uint16 {
	var num uint16
	binary.Read(b, binary.LittleEndian, &num)
	return num
}

func getHzOffset(b *bytes.Reader) int64 {
	b.Seek(4, 0)
	var mask byte
	binary.Read(b, binary.LittleEndian, &mask)
	if mask == 0x44 {
		return 0x2628
	} else if mask == 0x45 {
		return 0x26c4
	} else {
		fmt.Println("不支持的文件类型(无法获取汉语词组的偏移量)")
		os.Exit(1)
	}
	return -1
}

func getPyMap(b *bytes.Reader) map[uint16]string {
	pyMap := make(map[uint16]string)
	b.Seek(0x1540+4, 0)

	for {
		pyIdx := readUint16(b)
		pyLen := readUint16(b)
		pyStr := readUtf16Str(b, -1, int(pyLen))

		if _, ok := pyMap[pyIdx]; !ok {
			pyMap[pyIdx] = pyStr
		}

		if pyStr == "zuo" {
			break
		}
	}
	return pyMap
}

func getRecords(b *bytes.Reader, fileSize int64, hzOffset int64, pyMap map[uint16]string) []string {
	b.Seek(int64(hzOffset), io.SeekStart)
	var records []string
	for b.Size()-int64(b.Len()) != fileSize {
		wordCount := readUint16(b)
		pyIdxCount := int(readUint16(b) / 2)

		pySet := make([]string, pyIdxCount)
		for i := 0; i < pyIdxCount; i++ {
			pyIdx := readUint16(b)
			if py, ok := pyMap[pyIdx]; ok {
				pySet[i] = py
			} else {
				return records
			}
		}
		pyStr := strings.Join(pySet, "'")
		// 修正拼音
		pyStr = strings.Replace(pyStr, "lue", "lve", -1)
		pyStr = strings.Replace(pyStr, "nue", "nve", -1)

		for i := 0; i < int(wordCount); i++ {
			wordLen := readUint16(b)
			wordStr := readUtf16Str(b, -1, int(wordLen))

			// 跳过 ext_len 和 ext 共 12 个字节
			b.Seek(12, io.SeekCurrent)
			records = append(records, fmt.Sprintf("%s\t%s\t1", wordStr, pyStr))
		}
	}
	return records
}

func getWordsFromSogouCellDict(fname string) []string {
	data, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	b := bytes.NewReader(data)
	hzOffset := getHzOffset(b)
	pyMap := getPyMap(b)
	fileSize := int64(len(data))
	words := getRecords(b, fileSize, hzOffset, pyMap)

	return words
}

func save(records []string, f *os.File) []string {
	recordsTranslated := make([]string, len(records))
	for i, record := range records {
		recordsTranslated[i] = record
	}
	output := strings.Join(recordsTranslated, "\n")
	_, err := f.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
	return recordsTranslated
}
