package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
		t.Log("Success")
	default:
		t.Errorf("Type is not http.Handler but instead it is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
		t.Log("Success")
	default:
		t.Errorf("Type is not http.Handler but instead it is %T", v)
	}
}
