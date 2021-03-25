package hack

import (
	"testing"
)

func TestString(t *testing.T) {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	if String(b)!="hello" {
		t.Errorf("String mismatch: have %v, want %v.",String(b),"hello")
	}
}

func TestByte(t *testing.T) {
	s := "hello"
	v:=[]byte{'h', 'e', 'l', 'l', 'o'}
	if String(b)!="hello" {
		t.Errorf("Byte mismatch: have %v, want %v.",Byte(s),v)
	}
}