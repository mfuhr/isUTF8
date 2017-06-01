package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
)

const MaxInt = int(^uint(0) >> 1)

//
// The Unicode Standard
// Version 9.0 - Core Specification
// Chapter 3 Conformance
// 3.9 Unicode Encoding Forms
// Table 3-7. Well-Formed UTF-8 Byte Sequences
// http://www.unicode.org/versions/Unicode9.0.0/ch03.pdf#G7404
//
func bufferIsUTF8(fd int, offset int64, length int, checkSize int) (bool, int64) {
	if checkSize > length {
		log.Fatalf("checkSize (%v) > length (%v)", checkSize, length)
	}

	idxStart := offset % int64(unix.Getpagesize())
	idxEnd := checkSize + int(idxStart)
	mapOffset := offset - idxStart
	mapLength := length + int(idxStart)
	//fmt.Printf("offset=%v, mapOffset=%v, idxStart=%v, idxEnd=%v, length=%v, mapLength=%v\n",
	//           offset, mapOffset, idxStart, idxEnd, length, mapLength)

	buf, err := unix.Mmap(fd, mapOffset, mapLength, unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		log.Fatalf("Mmap(%v, %v): %s", offset, length, err)
	}

	bufSize := len(buf)
	if bufSize != mapLength {
		log.Fatalf("bufSize (%v) != mapLength (%v)", bufSize, mapLength)
	}

	isUTF8 := true

	var i int = int(idxStart)

bufLoop:
	for i < idxEnd {
		switch {
		case buf[i] <= 0x7f: // U+0000..U+007F
			i++
		case buf[i] >= 0xc2 && buf[i] <= 0xdf && // U+0080..U+07FF
			i+1 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0xbf:
			i += 2
		case buf[i] == 0xe0 && // U+0800..U+0FFF
			i+2 < bufSize &&
			buf[i+1] >= 0xa0 && buf[i+1] <= 0xbf &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf:
			i += 3
		case buf[i] >= 0xe1 && buf[i] <= 0xec && // U+1000..U+CFFF
			i+2 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0xbf &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf:
			i += 3
		case buf[i] == 0xed && // U+D000..U+D7FF
			i+2 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0x9f &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf:
			i += 3
		case buf[i] >= 0xee && buf[i] <= 0xef && // U+E000..U+FFFF
			i+2 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0xbf &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf:
			i += 3
		case buf[i] == 0xf0 && // U+10000..U+3FFFF
			i+3 < bufSize &&
			buf[i+1] >= 0x90 && buf[i+1] <= 0xbf &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf &&
			buf[i+3] >= 0x80 && buf[i+3] <= 0xbf:
			i += 4
		case buf[i] >= 0xf1 && buf[i] <= 0xf3 && // U+40000..U+FFFFF
			i+3 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0xbf &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf &&
			buf[i+3] >= 0x80 && buf[i+3] <= 0xbf:
			i += 4
		case buf[i] == 0xf4 && // U+100000..U+10FFFF
			i+3 < bufSize &&
			buf[i+1] >= 0x80 && buf[i+1] <= 0x8f &&
			buf[i+2] >= 0x80 && buf[i+2] <= 0xbf &&
			buf[i+3] >= 0x80 && buf[i+3] <= 0xbf:
			i += 4
		default:
			isUTF8 = false
			break bufLoop
		}
	}

	err = unix.Munmap(buf)
	if err != nil {
		log.Fatalf("Munmap: %s", err)
	}

	return isUTF8, int64(i) - idxStart
}

func fileIsUTF8(fname string) bool {
	f, err := unix.Open(fname, unix.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Open: %s: %s", fname, err)
	}

	defer unix.Close(f)

	var sbuf unix.Stat_t
	err = unix.Fstat(f, &sbuf)
	if err != nil {
		log.Fatalf("Fstat: %s: %s", fname, err)
	} else if (sbuf.Mode & unix.S_IFREG) == 0 {
		log.Fatalf("%s: not a regular file", fname)
	}

	fileSize := sbuf.Size
	var mapOffset int64 = 0
	maxMapSize := MaxInt - unix.Getpagesize() // leave room for alignment
	maxCharSize := 4

	if maxMapSize < maxCharSize {
		log.Fatalf("maxMapSize (%v) < maxCharSize (%v)", maxMapSize, maxCharSize)
	}

	for mapOffset < fileSize {
		bytesLeft := fileSize - mapOffset

		var mapSize int
		var checkSize int

		if bytesLeft > int64(maxMapSize) {
			mapSize = maxMapSize
			checkSize = mapSize - maxCharSize
		} else {
			mapSize = int(bytesLeft)
			checkSize = mapSize
		}

		//fmt.Printf("fileSize=%v, mapOffset=%v, bytesLeft=%v, mapSize=%v, checkSize=%v\n",
		//            fileSize, mapOffset, bytesLeft, mapSize, checkSize)

		isUTF8, bytesChecked := bufferIsUTF8(f, mapOffset, mapSize, checkSize)
		//fmt.Printf("isUTF8=%v, bytesChecked=%v\n", isUTF8, bytesChecked)

		if !isUTF8 {
			return false
		}

		mapOffset += bytesChecked
		//fmt.Printf("mapOffset=%v\n\n", mapOffset)
	}

	return true
}

func main() {
	exitCode := 0

	if len(os.Args) > 2 {
		log.Printf("warning: multiple filename arguments\n")
	}

	for _, arg := range os.Args[1:] {
		isUTF8 := fileIsUTF8(arg)
		fmt.Printf("%v %s\n", isUTF8, arg)
		if !isUTF8 {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
