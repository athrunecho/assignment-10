package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	responseHeader := w.Header()
	for k, v := range r.Header {
		responseHeader.Add(k, fmt.Sprint(v))
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	fmt.Fprintf(w, "Hello")
	log.Printf("client IP: %v", getCurrentIP(r))
	log.Printf("client response code: %v", 200)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "200")
}

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	timer := NewTimer()
	defer timer.ObserveTotal()
	randInt := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	w.Write([]byte(fmt.Sprintf("<h1>%d<h1>", randInt)))
}

func main() {
	Register()
	os.Setenv("VERSION", "1.0")
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/healthz", HealthHandler)
	http.HandleFunc("/images", ImagesHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
