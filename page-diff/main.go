package main

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/n7olkachev/imgdiff/pkg/imgdiff"
)

const (
	css = `
	.pub{
		display: none !important;
	}
	.contentVisibility{
		content-visibility: visible !important;
	}`
	addCssScript = `
	(css) => {
		const style = document.createElement('style');
		style.type = 'text/css';
		style.appendChild(document.createTextNode(css));
		document.head.appendChild(style);
		return true;
	}
	`
)

func generateScreenshot(ctx context.Context, url, filename string, headers map[string]any) (image.Image, error) {
	var data []byte
	tasks := chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(headers),
		chromedp.Navigate(url),
		chromedp.PollFunction(addCssScript, nil, chromedp.WithPollingArgs(css)),
		chromedp.Evaluate("window.scrollTo(0, document.body.scrollHeight);", nil),
		chromedp.Sleep(4 * time.Second),
		chromedp.Evaluate("window.scrollTo(0, 0);", nil),
		chromedp.FullScreenshot(&data, 90),
	}
	if err := chromedp.Run(ctx, tasks); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func main() {
	// Define chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WindowSize(1920, 1080), chromedp.Flag("headless", false))

	// Create chrome context
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Get reference image
	referenceImg, err := generateScreenshot(ctx, "https://www.google.com/", "./reference.jpeg", map[string]any{})
	if err != nil {
		log.Panic(err)
	}

	// Get new image
	newImg, err := generateScreenshot(ctx, "https://www.example.com", "./new.jpeg", map[string]any{})
	if err != nil {
		log.Panic(err)
	}

	// Compare images
	result := imgdiff.Diff(referenceImg, newImg, &imgdiff.Options{
		Threshold: 0.1,
		DiffImage: false,
	})

	diffImg, err := os.Create("./diff.png")
	if err != nil {
		log.Panic(err)
	}
	defer diffImg.Close()

	err = png.Encode(diffImg, result.Image)
	if err != nil {
		log.Panic(err)
	}
}
