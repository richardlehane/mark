package main

import (
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	wm := pdfcpu.DefaultWatermarkConfig()
	wm.TextLines = []string{"Richard", "Rules"}
	log.Fatal(api.AddWatermarksFile("test.pdf", "", nil, wm, nil))
}
