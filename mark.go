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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

var (
	overwritef = flag.Bool("d", false, "delete original files after watermarking")
	statsf     = flag.Bool("s", false, "print a statistics files with details about the PDFs watermarked")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Expects a single argument with the name of a pdf file or directory")
		os.Exit(1)
	}
	target := flag.Arg(0)
	// get researcher's name from stdin
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
	// if -s flag used, keep stats
	var statsdir string
	var pdfnames, errortxts []string
	var totalPages int
	var count int
	if *statsf {
		pdfnames = make([]string, 0, 500)
	}
	// walk the target dir
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking %q: %v", path, err)
		}
		if info.IsDir() {
			if *statsf && target == path {
				statsdir = target
			}
			return nil
		}
		if filepath.Ext(path) == ".pdf" {
			// don't watermark files that are already watermarked
			if strings.HasSuffix(path, "_wm.pdf") {
				return nil
			}
			// if target == path, this is a single PDF watermarking
			if *statsf && target == path {
				statsdir = filepath.Dir(target)
			}
			ctx, err := pdfcpu.ReadFile(path, pdfcpu.NewDefaultConfiguration())
			if err == nil {
				err = validate.XRefTable(ctx.XRefTable)
			}
			if err == nil {
				err = pdfcpu.OptimizeXRefTable(ctx)
			}
			if err == nil {
				err = ctx.EnsurePageCount()
			}
			var pageCount int
			if err == nil {
				// make a new watermark
				wm := pdfcpu.DefaultWatermarkConfig()
				wm.Opacity = 0.5
				wm.TextLines = text
				// select all pages
				pageCount = ctx.PageCount
				m := pdfcpu.IntSet{}
				for i := 1; i <= pageCount; i++ {
					m[i] = true
				}
				// add watermark
				err = pdfcpu.AddWatermarks(ctx, m, wm)
			}
			// validate again
			if err == nil {
				err = validate.XRefTable(ctx.XRefTable)
			}
			// now write to file
			if err == nil {
				nn := strings.TrimSuffix(path, ".pdf") + "_wm.pdf"
				nf, err := os.Create(nn)
				if err == nil {
					err = api.WriteContext(ctx, nf)
				}
				nf.Close()
			}
			// if we've successfully written a watermarked file, update stats
			if err == nil && *statsf {
				count++
				totalPages += pageCount
				if *statsf {
					pdfnames = append(pdfnames, filepath.Base(path))
				}
			}
			// if we've successfully written a watermaked file, delete the old if in overwrite mode
			if err == nil && *overwritef {
				err = os.Remove(path)
			}
			// if errors in any of above: write to stats and clear the error if we are in stats mode; otherwise escalate the error
			if err != nil {
				if *statsf {
					errortxts = append(errortxts, fmt.Sprintf("error for file %q: %v", path, err))
					return nil
				}
				return err
			}
		}
		return nil
	})
	if statsdir != "" && count > 0 {
		sf, err := os.Create(filepath.Join(statsdir, fmt.Sprintf("stats-%d.txt", time.Now().Unix())))
		if err == nil {
			fmt.Fprintf(sf,
				"WATERMARK STATS\n---\nResearcher: %s\nDate: %s\nTotal files: %d\nTotal pages: %d\n---\nFiles:\n%s",
				name,
				text[1],
				count,
				totalPages,
				strings.Join(pdfnames, "\n"),
			)
			if len(errortxts) > 0 {
				fmt.Fprintf(sf, "\n---\nErrors:\n%s", strings.Join(errortxts, "\n"))
			}
		}
		sf.Close()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
