// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package webapp handles the web application.
package webapp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twihike/go-geojp/pkg/geo/jp"
	"github.com/twihike/go-structconv/structconv"
)

type appConfig struct {
	Port           string
	StaticDir      string
	StaticURL      string
	HealthCheckURL string
	AddrPosPath    string
}

var (
	conf appConfig = appConfig{
		Port:           "8080",
		AddrPosPath:    "latest.csv",
		HealthCheckURL: "/api/health",
		StaticDir:      "web/static",
		StaticURL:      "/",
	}
	aps  jp.AddressPositions
	iaps jp.IndexedAPs
)

// RunServer runs the web application server.
func RunServer() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("starting app...")

	if err := structconv.DecodeEnv(&conf); err != nil {
		log.Fatal(err)
	}
	log.Printf("config: %+v\n", conf)

	var err error
	aps, err = jp.ReadAPsFromFile(conf.AddrPosPath)
	if err != nil {
		log.Fatalln(err)
	}
	iaps = jp.CreateIndexedAPs(aps)

	server := setupServer()
	runServer(server)
}

func setupServer() *http.Server {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(conf.StaticDir))
	mux.Handle(conf.StaticURL, http.StripPrefix(conf.StaticURL, files))
	mux.HandleFunc(conf.HealthCheckURL, health)
	mux.HandleFunc("/api/geocoding", geocoding)
	mux.HandleFunc("/api/reverse-geocoding", reverseGeocoding)
	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}
	return server
}

func runServer(server *http.Server) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		log.Println("stopping app...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Println(err)
		}
		close(idleConnsClosed)
	}()

	log.Println("app started on port:", conf.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-idleConnsClosed
	log.Println("app stopped")
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b := map[string]string{"status": "up"}
	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
