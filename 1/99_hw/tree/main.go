package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

const (
	ISHape    = "│"
	TShape    = "├"
	LShape    = "└"
	DashShape = "─"
)

func tabify(output io.Writer) io.Writer {
	fmt.Fprint(output, "\t")
	return output
}

func entryInfo(e fs.DirEntry) (string, error) {
	fi, err := e.Info()
	if err != nil {
		return "", err
	}

	infoString := fi.Name()
	if !fi.IsDir() {
		s := fi.Size()
		if s == 0 {
			infoString += " " + "(empty)"
		} else {
			infoString += " " + fmt.Sprintf("(%db)", s)
		}
	}
	return infoString, nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() && !printFiles {
			continue
		}
		info, err := entryInfo(e)
		if err != nil {
			return err
		}
		println(info)
		if e.IsDir() {
			err = dirTree(tabify(output), path+"/"+e.Name(), printFiles)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
