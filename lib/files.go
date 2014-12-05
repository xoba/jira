package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	var list fileList
	filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) (out error) {
		if strings.HasSuffix(path, ".nn.go") {
			return
		}
		if strings.HasSuffix(path, ".y.go") {
			return
		}
		if strings.Contains(path, "flymake_") && strings.HasSuffix(path, ".go") {
			return
		}
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".nex") || strings.HasSuffix(path, ".y") {
			list = append(list, file{path, info.ModTime()})
		}
		return
	})
	sort.Sort(list)
	for _, f := range list {
		fmt.Println(f.name)
	}
}

type file struct {
	name string
	time time.Time
}

type fileList []file

func (f fileList) Len() int {
	return len(f)
}
func (f fileList) Less(i, j int) bool {
	if f[i].time == f[j].time {
		return f[i].name < f[j].name
	} else {
		return f[i].time.Before(f[j].time)
	}
}
func (f fileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
