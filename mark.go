package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

var (
	overwrite = flag.Bool("d", false, "delete original files after watermarking")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Expects a single argument with the name of a pdf file or directory")
		os.Exit(1)
	}
	target := flag.Arg(0)
	// get watermark name from stdin
	reader := bufio.NewReader(os.Stdin)
	var name string
L:
	for {
		fmt.Print("Enter name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		fmt.Printf("Watermark with '%s'? (y)es/ enter, (e)dit, (q)uit: ", name)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		switch cmd {
		case "", "y", "yes":
			break L
		case "q", "quit":
			fmt.Println("quitting")
			os.Exit(0)
		}
	}
	text := []string{
		name,
		fmt.Sprintf("made available on %s", time.Now().Format("2006-01-02")),
	}
	// walk the target dir
	var count int
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking %q: %v\n", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".pdf" {
			// don't watermark files that are already watermarked
			if strings.HasSuffix(path, "_wm.pdf") {
				return nil
			}
			count++
			wm := pdfcpu.DefaultWatermarkConfig()
			wm.Opacity = 0.5
			wm.TextLines = text
			nn := strings.TrimSuffix(path, ".pdf") + "_wm.pdf"
			if *overwrite {
				err := os.Rename(path, nn)
				if err != nil {
					return err
				}
				return api.AddWatermarksFile(nn, "", nil, wm, nil)
			}
			return api.AddWatermarksFile(path, nn, nil, wm, nil)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("finished! %d pdf files watermarked\n", count)
}
