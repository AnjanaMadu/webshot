package main

import (
	"context"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte(`{"status": "ok"}`))
		}
	})
	http.HandleFunc("/api", apiHandler)
	log.Printf("[INFO] Listening on port 9090")
	http.ListenAndServe(":9090", nil)
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
		w.Write([]byte(`{"status": "error", "message": "URL is required"}`))
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
