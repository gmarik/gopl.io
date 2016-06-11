package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Exercise 5.14: Use the breadthFirst function to explore a different structure. For example, you could use the course dependencies from the topoSort example (a directed graph), the file system hierarchy on your computer (a tree), or a list of bus or subway routes downloaded from your city government’s web site (an undirected graph)

func main() {
	err := breadthFirst(walk(func(s string) { log.Println(s) }), os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func walk(f func(string)) func(string) ([]string, error) {
	return func(path string) ([]string, error) {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		f(path)

		if !info.IsDir() {
			return nil, nil
		}

		fileInfos, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}

		paths := make([]string, len(fileInfos))

		for i := 0; i < len(fileInfos); i += 1 {
			info := fileInfos[i]
			p := filepath.Join(path, info.Name())
			f(p)
			paths[i] = p
		}

		return paths, nil
	}

}

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) ([]string, error), worklist []string) error {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				childItems, err := f(item)
				if err != nil {
					return err
				}
				worklist = append(worklist, childItems...)
			}
		}
	}

	return nil
}
