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
	return dirInTree(out, path, printFiles, 0)
}

func dirInTree(out *os.File, path string, printFiles bool, depth int) error {
	de, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	sortDirs(de)
	for _, d := range de {
		if d.IsDir() {
			dirPrint(out, d.Name(), depth)

			err := dirInTree(out, path+"/"+d.Name(), printFiles, depth+1)
			if err != nil {
				return nil
			}
		} else {
			if printFiles {
				dirPrint(out, d.Name(), depth)
			}
		}
	}
	return nil
}

const (
	begin = iota
	mid
	end
)

var syms = map[int]string{
	begin: "├───",
	mid: "│	",
	end: "└───",
}

func sortDirs(dd []os.DirEntry) {
	sort.Slice(dd, func(i, j int) bool { return dd[i].Name() < dd[j].Name() })
}

func dirPrint(out *os.File, name string, depth int) {
	// if i == 0 {
	// 	fmt.Fprintf(out, "%s%s/n", syms[begin], name)
	// } else if i < len {
	// 	fmt.Fprintf(out, "%s%s/n", syms[begin], name)
	// }
	for ; depth > 0; depth-- {
		fmt.Fprintf(out, "\t")
	}
	fmt.Fprintf(out, "%s\n", name)
}

func genPrefix(prefix string, i, len int, isInLast bool) string {
	if isInLast {
		prefix += "\t"
	} else {
		prefix += syms[mid] + "\t"
	}

	if i < len {
		prefix += syms[begin]
	} else {
		prefix += syms[end]
	}

	return prefix
}
