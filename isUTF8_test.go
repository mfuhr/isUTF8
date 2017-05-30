package main

import (
	"path/filepath"
	"testing"
)

const testDir = "testdata"

func testFile(t *testing.T, fileName string, expectedResult bool) {
	if fileIsUTF8(filepath.Join(testDir, fileName)) != expectedResult {
		t.Fail()
	}
}

func TestUTF8(t *testing.T) {
	testFile(t, "test_utf8.txt", true)
}

func TestLatin1(t *testing.T) {
	testFile(t, "test_latin1.txt", false)
}

func TestASCII(t *testing.T) {
	testFile(t, "test_ascii.txt", true)
}

func TestAllUTF8(t *testing.T) {
	testFile(t, "test_all_utf8.txt", true)
}

func TestUTF8Short2ByteC2(t *testing.T) {
	testFile(t, "test_utf8_short_2byte_c2.txt", false)
}

func TestUTF8Short3ByteE0(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_e0.txt", false)
}

func TestUTF8Short3ByteE1(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_e1.txt", false)
}

func TestUTF8Short3ByteED(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_ed.txt", false)
}

func TestUTF8Short3ByteEE(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_ee.txt", false)
}

func TestUTF8Short4ByteF0(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f0.txt", false)
}

func TestUTF8Short4ByteF1(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f1.txt", false)
}

func TestUTF8Short4ByteF4(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f4.txt", false)
}
