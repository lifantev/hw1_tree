package main

import (
	// "bytes"
	"fmt"
	"sort"

	// "io"
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

func dirTree(out *os.File, path string, printFiles bool) (err error) {
	return dirInTree(out, path, "", printFiles)
}

func dirInTree(out *os.File, path, indent string, printFiles bool) error {
	de, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var ln int
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
	sortDirs(de)

	for i, d := range de {
		if d.IsDir() {
			dirPrint(out, indent, genPrefix(i, ln), d.Name())

			err := dirInTree(out, path+"/"+d.Name(), genIndent(indent, i, ln), printFiles)
			if err != nil {
				return nil
			}
		} else {
			if printFiles {
				dirPrint(out, indent, genPrefix(i, ln), d.Name())
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

func sortDirs(dd []os.DirEntry) {
	sort.Slice(dd, func(i, j int) bool { return dd[i].Name() < dd[j].Name() })
}

func dirPrint(out *os.File, indent, prefix, name string) {
	fmt.Fprintf(out, "%s%s%s\n", indent, prefix, name)
}

func genIndent(indent string, i, len int) string {
	if i < len-1 {
		indent += syms[mid]
	} else {
		indent += syms[tab]
	}
	return indent
}

func genPrefix(i, len int) string {
	if i < len-1 {
		return syms[begin]
	} else {
		return syms[end]
	}
}

// todo:  suffix of file length
