package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed staticAssets
var staticAssets embed.FS

type BocaMockServer struct {
	listenAddress string
}

func NewBocaMockServer(listenAddress string) *BocaMockServer {
	return &BocaMockServer{listenAddress: listenAddress}
}

func (b *BocaMockServer) Start() error {
	router := mux.NewRouter().StrictSlash(true)
	dir, err := fs.Sub(staticAssets, "staticAssets")
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(http.FileServer(http.FS(dir)))
	return http.ListenAndServe(b.listenAddress, router)
}
