package main

import (
	"context"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func main() {
	server := http.Server{Addr: ":9090"}
	server.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				sendJSON(w, 200, `{"status": "ok", "message": "Welcome to screenshot API"}`)
			case "/api":
				apiHandler(w, r)
			default:
				sendJSON(w, 404, `{"status": "error", "message": "Not found"}`)
			}
		},
	)
	log.Printf("[INFO] Listening on port 9090")
	server.ListenAndServe()
}

func sendJSON(w http.ResponseWriter, status int, data string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(data))
}

var chromeOpts = []func(allocator *chromedp.ExecAllocator){
	chromedp.ExecPath("/snap/bin/chromium"),
	chromedp.Flag("disable-dev-shm-usage", true),
	chromedp.Flag("disable-background-networking", true),
	chromedp.WindowSize(1280, 720),
	chromedp.DisableGPU,
	chromedp.NoFirstRun,
	chromedp.NoSandbox,
	chromedp.Headless,
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		return
	}
	link := r.URL.Query().Get("url")
	if link == "" {
		sendJSON(w, 400, `{"status": "error", "message": "URL is required"}`)
		return
	}
	log.Printf("[INFO] Link: %s", link)

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(
		context.Background(),
		chromeOpts...,
	)
	defer allocatorCancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	var buf []byte
	var title string
	chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(link),
		chromedp.Title(&title),
		chromedp.FullScreenshot(&buf, 90),
	})

	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("X-Page-Title", title)
	w.Write(buf)
}
