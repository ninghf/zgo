package container

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := new(ConcurrentMap)

	if v := m.Get("key"); v!= nil {
		t.Errorf("get mismatch: have %v, want %v.",v,nil)
	}

	setv:= "value";
	m.Set("key",setv)
	if v := m.Get("key"); v!= setv {
		t.Errorf("set fail, get mismatch: have %v, want %v.",v,setv)
	}

	m.Del("key")
	if v := m.Get("key"); v!= nil {
		t.Errorf("del fail, get mismatch: have %v, want %v.",v,nil)
	}

	m.Set(1, "1")
	m.Set(2, 2)
	m.Set("3", 3)

	if v := m.Len(); v!= 3 {
		t.Errorf("get len fail,mismatch: have %v, want %v.",v,3)
	}
}