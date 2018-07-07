package SimpleHttpServer

import (
	"net/http"
)

type SimpleHttpServer struct {
	httpServer *http.Server
}

func New_SimpleHttpServer(listenAddr string) *SimpleHttpServer {
	var curData *SimpleHttpServer = new(SimpleHttpServer)
	//
	curData.httpServer = new(http.Server)
	curData.httpServer.Addr = listenAddr
	curData.httpServer.Handler = http.NewServeMux()
	//
	return curData
}

func (self *SimpleHttpServer) GetHttpServeMux() *http.ServeMux {
	return self.httpServer.Handler.(*http.ServeMux)
}

func (self *SimpleHttpServer) Run() error {
	return self.httpServer.ListenAndServe()
}

func (self *SimpleHttpServer) RunTLS(certFile string, keyFile string) error {
	return self.httpServer.ListenAndServeTLS(certFile, keyFile)
}
