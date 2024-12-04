package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartServer(port int) {
	addr := ":" + strconv.Itoa(port)
	log.Printf("Starting server listening at %s", addr)

	server := &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     NewAuthHandlerBuilder().Build(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
