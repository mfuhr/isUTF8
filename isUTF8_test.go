package main

import "testing"

func TestUTF8(t *testing.T) {
  if isUTF8 := fileIsUTF8("test_utf8.txt"); !isUTF8 {
    t.Fail()
  }
}

func TestLatin1(t *testing.T) {
  if isUTF8 := fileIsUTF8("test_latin1.txt"); isUTF8 {
    t.Fail()
  }
}
