package main

import (
	// "bytes"
	"fmt"
	"sort"

	"io"
	"os"
	// "path/filepath"
	// "strings"
)

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

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	return dirInTree(out, path, "", printFiles)
}

func dirInTree(out io.Writer, path, indent string, printFiles bool) error {
	de, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var ln int
	de, ln = filterDirs(printFiles, de, ln)
	sortDirs(de)

	for i, d := range de {
		isLast := i == (ln - 1)
		relPath := path + "/" + d.Name()
		s, err := getSuffix(relPath)
		if err != nil {
			return err
		}

		if d.IsDir() {
			dirPrint(out, indent, genPrefix(isLast), d.Name(), "")
			err = dirInTree(out, relPath, genIndent(indent, isLast), printFiles)
			if err != nil {
				return nil
			}
		} else {
			if printFiles {
				dirPrint(out, indent, genPrefix(isLast), d.Name(), s)
			}
		}
	}
	return nil
}

const (
	begin = iota
	tab
	mid
	end
)

var syms = map[int]string{
	begin: "├───",
	tab: "	",
	mid: "│	",
	end: "└───",
}

func filterDirs(printFiles bool, de []os.DirEntry, ln int) ([]os.DirEntry, int) {
	if !printFiles {
		n := 0
		for _, d := range de {
			if d.IsDir() {
				de[n] = d
				n++
			}
		}
		de = de[:n]
		ln = n
	} else {
		ln = len(de)
	}
	return de, ln
}

func sortDirs(dd []os.DirEntry) {
	sort.Slice(dd, func(i, j int) bool { return dd[i].Name() < dd[j].Name() })
}

func dirPrint(out io.Writer, indent, prefix, name, suffix string) {
	fmt.Fprintf(out, "%s%s%s%s\n", indent, prefix, name, suffix)
}

func genIndent(indent string, isLast bool) string {
	if !isLast {
		indent += syms[mid]
	} else {
		indent += syms[tab]
	}
	return indent
}

func genPrefix(isLast bool) string {
	if !isLast {
		return syms[begin]
	} else {
		return syms[end]
	}
}

func getSuffix(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}

	size := fi.Size()
	if size > 0 {
		return fmt.Sprintf(" (%db)", size), nil
	} else {
		return " (empty)", nil
	}
}
