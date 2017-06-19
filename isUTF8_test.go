package main

import (
	"golang.org/x/sys/unix"
	"path/filepath"
	"testing"
)

const testDir = "testdata"

func testFile(t *testing.T, fileName string, expectedResult bool, errorExpected bool) {
	testMaxInt := 2 * unix.Getpagesize()

	isUTF8, err := fileIsUTF8(filepath.Join(testDir, fileName), testMaxInt)
	if err != nil && !errorExpected {
		t.Fail()
	} else if isUTF8 != expectedResult {
		t.Fail()
	}
}

func TestUTF8(t *testing.T) {
	testFile(t, "test_utf8.txt", true, false)
}

func TestLatin1(t *testing.T) {
	testFile(t, "test_latin1.txt", false, false)
}

func TestASCII(t *testing.T) {
	testFile(t, "test_ascii.txt", true, false)
}

func TestAllUTF8(t *testing.T) {
	testFile(t, "test_all_utf8.txt", true, false)
}

func TestUTF8Short2ByteC2(t *testing.T) {
	testFile(t, "test_utf8_short_2byte_c2.txt", false, false)
}

func TestUTF8Short3ByteE0(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_e0.txt", false, false)
}

func TestUTF8Short3ByteE1(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_e1.txt", false, false)
}

func TestUTF8Short3ByteED(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_ed.txt", false, false)
}

func TestUTF8Short3ByteEE(t *testing.T) {
	testFile(t, "test_utf8_short_3byte_ee.txt", false, false)
}

func TestUTF8Short4ByteF0(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f0.txt", false, false)
}

func TestUTF8Short4ByteF1(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f1.txt", false, false)
}

func TestUTF8Short4ByteF4(t *testing.T) {
	testFile(t, "test_utf8_short_4byte_f4.txt", false, false)
}

func TestUTF8SurrogateD800(t *testing.T) {
	testFile(t, "test_utf8_surrogate_d800.txt", false, false)
}

func TestUTF8SurrogateDB7F(t *testing.T) {
	testFile(t, "test_utf8_surrogate_db7f.txt", false, false)
}

func TestUTF8SurrogateDB80(t *testing.T) {
	testFile(t, "test_utf8_surrogate_db80.txt", false, false)
}

func TestUTF8SurrogateDBFF(t *testing.T) {
	testFile(t, "test_utf8_surrogate_dbff.txt", false, false)
}

func TestUTF8SurrogateDC00(t *testing.T) {
	testFile(t, "test_utf8_surrogate_dc00.txt", false, false)
}

func TestUTF8SurrogateDFFF(t *testing.T) {
	testFile(t, "test_utf8_surrogate_dfff.txt", false, false)
}

func TestDirectory(t *testing.T) {
	testFile(t, ".", false, true)
}

func TestMissingFile(t *testing.T) {
	testFile(t, "no_such_file.txt", false, true)
}

func TestCheckSize(t *testing.T) {
	_, _, err := bufferIsUTF8(0, 0, 1, 2)
	if err == nil {
		t.Fail()
	}
}

func TestMmapFailure(t *testing.T) {
	_, _, err := bufferIsUTF8(0, 0, 1, 1)
	if err == nil {
		t.Fail()
	}
}

func TestMaxIntTooSmall(t *testing.T) {
	_, err := fileIsUTF8(filepath.Join(testDir, "test_utf8.txt"), unix.Getpagesize() - 1)
	if err == nil {
		t.Fail()
	}
}
