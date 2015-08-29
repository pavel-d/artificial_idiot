package tools

import (
	"fmt"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Health check request: %s %s", req.RemoteAddr, req.RequestURI)
}

func StartHealthCheck(port int) {
	http.HandleFunc("/", HelloServer)
	log.Printf("Starting health check handler on %v port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
