package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-chi-cache/cache"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong " + fmt.Sprint(time.Now().Unix())))
	})

	cached := cache.NewCache(false)
	cached.SetPaths([]string{"/ping_cache_good"})

	r.With(cached.CachePage()).Get("/ping_cache", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong " + fmt.Sprint(time.Now().Unix())))
	})
	// 指定路径
	r.With(cached.CachePage()).Get("/ping_cache_good", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong " + fmt.Sprint(time.Now().Unix())))
	})

	// 授权登陆
	cached2 := cache.NewCache(true)
	r.With(cached2.CachePage()).Get("/auth_test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong " + fmt.Sprint(time.Now().Unix())))
	})

	cached3 := cache.NewCache(true)
	r.With(cached3.CachePageWithQuery()).Get("/auth_test_query", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong " + fmt.Sprint(time.Now().Unix())))
	})

	log.Println("open: http://127.0.0.1:9652/ping")
	log.Println("open: http://127.0.0.1:9652/ping_cache")

	http.ListenAndServe(":9652", r)
}
