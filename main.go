package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/chromedp/chromedp"
)

func main() {
	r := gin.Default()
	r.GET("/", indexHandler)
	r.GET("/api", apiHandler)
	r.Run(":9090")
}

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
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

func apiHandler(c *gin.Context) {
	link := c.Query("url")
	if link == "" {
		c.JSON(400, gin.H{"status": "error", "message": "URL is required"})
		return
	}

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

	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("X-Page-Title", title)
	c.Data(200, "image/jpeg", buf)
}
