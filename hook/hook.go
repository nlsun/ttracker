package hook

import (
	"net"
	"net/http"

	"github.com/nlsun/ttracker/output"
)

type hook struct {
	lsnr net.Listener
}

type Server interface {
	Serve(output.Logger) error
	Addr() string
}

func NewServer(lsnr net.Listener) Server {
	return hook{lsnr: lsnr}
}

func (s hook) Serve(l output.Logger) error {
	mux := http.NewServeMux()
	mux.Handle("/hook", l)
	return http.Serve(s.lsnr, mux)
}

func (s hook) Addr() string {
	return s.lsnr.Addr().String()
}
