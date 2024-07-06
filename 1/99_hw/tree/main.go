package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

const (
	ISHape    = "│"
	TShape    = "├"
	LShape    = "└"
	DashShape = "─"
)

func countEntries(entries []fs.DirEntry, printFiles bool) (count int) {
	for _, e := range entries {
		// Skip files.
		if !e.IsDir() && !printFiles {
			continue
		}
		count++
	}
	return
}

func entryInfo(e fs.DirEntry) (info string, err error) {
	fi, err := e.Info()
	if err != nil {
		return
	}

	info = fi.Name()
	if !fi.IsDir() {
		s := fi.Size()
		if s == 0 {
			info += " " + "(empty)"
		} else {
			info += " " + fmt.Sprintf("(%db)", s)
		}
	}
	return
}

func dirTreePrefixed(output io.Writer, path string, printFiles bool, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	entriesCount := countEntries(entries, printFiles)
	idx := 0

	for _, e := range entries {
		// Skip files.
		if !e.IsDir() && !printFiles {
			continue
		}
		info, err := entryInfo(e)
		if err != nil {
			return err
		}

		isLastEntry := entriesCount == 1 || idx+1 == entriesCount
		if isLastEntry {
			fmt.Fprint(output, prefix+LShape+strings.Repeat(DashShape, 3)+info+"\n")
		} else {
			fmt.Fprint(output, prefix+TShape+strings.Repeat(DashShape, 3)+info+"\n")
		}

		if e.IsDir() {
			var err error
			if isLastEntry {
				err = dirTreePrefixed(output, path+"/"+e.Name(), printFiles, prefix+"\t")
			} else {
				err = dirTreePrefixed(output, path+"/"+e.Name(), printFiles, prefix+ISHape+"\t")
			}

			if err != nil {
				return err
			}
		}
		idx++
	}

	return nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	return dirTreePrefixed(output, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
