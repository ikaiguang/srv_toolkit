package tke

import (
	"testing"
)

func TestFromError(t *testing.T) {
	code := Unknown
	err := New(code)
	s, ok := FromError(err)
	if !ok {
		t.Log("cannot parse error")
		t.Fail()
		return
	}
	if s.Code32() != code.Code() {
		t.Log("parse error success, but s.Code not equal err.Code")
		t.Fail()
		return
	}
}
