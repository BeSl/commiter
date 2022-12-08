package server

import (
	"net/http"
)

func (ms *MServer) createGatewayServer(host_port string) *http.Server {

	gatewayServer := &http.Server{
		Addr:    host_port,
		Handler: ms.newMuxServe(),
	}

	return gatewayServer
}

func (ms *MServer) newMuxServe() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/uploadtoquery", ms.uploadtoquery)
	mux.HandleFunc("/ping", ms.pingService)
	mux.HandleFunc("/crtab", ms.CreateTables)

	return mux
}
