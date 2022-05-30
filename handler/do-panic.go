package handler

import (
	"net/http"
)

type PanicHandler struct {
}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Panic!!!!!!!!!!!")
}
