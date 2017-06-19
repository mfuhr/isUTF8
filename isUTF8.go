package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
)

const sysMaxInt = int(^uint(0) >> 1)
const maxCharSize = 4

//
// The Unicode Standard
// Version 9.0 - Core Specification
// Chapter 3 Conformance
// 3.9 Unicode Encoding Forms
// Table 3-7. Well-Formed UTF-8 Byte Sequences
// http://www.unicode.org/versions/Unicode9.0.0/ch03.pdf#G7404
//
func bufferIsUTF8(fd int, offset int64, length int, checkSize int) (isUTF8 bool, bytesChecked int64, err error) {
	if checkSize > length {
		return false, 0, fmt.Errorf("checkSize (%v) > length (%v)", checkSize, length)
	}

	idxStart := offset % int64(unix.Getpagesize())
	idxEnd := checkSize + int(idxStart)
	mapOffset := offset - idxStart
	mapLength := length + int(idxStart)
	//fmt.Printf("offset=%v, mapOffset=%v, idxStart=%v, idxEnd=%v, length=%v, mapLength=%v\n",
	//           offset, mapOffset, idxStart, idxEnd, length, mapLength)

	buf, err := unix.Mmap(fd, mapOffset, mapLength, unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return false, 0, fmt.Errorf("Mmap(%v, %v): %s", offset, length, err)
	}

	defer func() {
		if localErr := unix.Munmap(buf); localErr != nil {
			fmt.Fprintf(os.Stderr, "Munmap error: %v\n", localErr)
			err = localErr
		}
	}()

	bufSize := len(buf)
	if bufSize != mapLength {
		return false, 0, fmt.Errorf("bufSize (%v) != mapLength (%v)", bufSize, mapLength)
	}

	isUTF8 = true

	i := int(idxStart)

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

	return isUTF8, int64(i) - idxStart, nil
}

func fileIsUTF8(fileName string, maxInt int) (isUTF8 bool, err error) {
	f, err := unix.Open(fileName, unix.O_RDONLY, 0)
	if err != nil {
		return false, err
	}

	defer func() {
		if localErr := unix.Close(f); localErr != nil {
			err = localErr
		}
	}()

	var sbuf unix.Stat_t
	err = unix.Fstat(f, &sbuf)
	if err != nil {
		return false, err
	} else if (sbuf.Mode & unix.S_IFREG) == 0 {
		return false, fmt.Errorf("%s: not a regular file", fileName)
	}

	fileSize := sbuf.Size

	maxMapSize := maxInt - unix.Getpagesize() // leave room for alignment
	if maxMapSize < maxCharSize {
		return false, fmt.Errorf("maxMapSize (%v) < maxCharSize (%v)", maxMapSize, maxCharSize)
	}

	for mapOffset := int64(0); mapOffset < fileSize; {
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

		isUTF8, bytesChecked, err := bufferIsUTF8(f, mapOffset, mapSize, checkSize)
		//fmt.Printf("isUTF8=%v, bytesChecked=%v, err=%v\n", isUTF8, bytesChecked, err)
		if err != nil {
			return false, err
		} else if !isUTF8 {
			return false, nil
		}

		mapOffset += bytesChecked
		//fmt.Printf("mapOffset=%v\n\n", mapOffset)
	}

	return true, nil
}

func main() {
	var maxInt int

	flag.IntVar(&maxInt, "-maxint", sysMaxInt, "maximum integer size (intended only for testing)")
	flag.Parse()

	exitCode := 0

	if len(os.Args) > 2 {
		log.Printf("warning: multiple filename arguments\n")
	}

	for _, arg := range os.Args[1:] {
		isUTF8, err := fileIsUTF8(arg, maxInt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %s\n", isUTF8, arg)
		if !isUTF8 {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
