package main

import (
  "path/filepath"
  "testing"
)

const testDir = "testdata"

func TestUTF8(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8.txt")); !isUTF8 {
    t.Fail()
  }
}

func TestLatin1(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_latin1.txt")); isUTF8 {
    t.Fail()
  }
}

func TestASCII(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_ascii.txt")); !isUTF8 {
    t.Fail()
  }
}

func TestAllUTF8(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_all_utf8.txt")); !isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short2ByteC2(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_2byte_c2.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short3ByteE0(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_3byte_e0.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short3ByteE1(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_3byte_e1.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short3ByteED(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_3byte_ed.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short3ByteEE(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_3byte_ee.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short4ByteF0(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_4byte_f0.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short4ByteF1(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_4byte_f1.txt")); isUTF8 {
    t.Fail()
  }
}

func TestUTF8Short4ByteF4(t *testing.T) {
  if isUTF8 := fileIsUTF8(filepath.Join(testDir, "test_utf8_short_4byte_f4.txt")); isUTF8 {
    t.Fail()
  }
}

