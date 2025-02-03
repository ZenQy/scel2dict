package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Dict struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	var dicts []Dict
	data, err := os.ReadFile("dict.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &dicts)
	if err != nil {
		fmt.Println("Error Unmarshal json file:", err)
		os.Exit(1)
	}

	dirs := []string{"out", "scel"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}

	for _, dict := range dicts {
		dict.Download()
	}

	scelFiles, _ := filepath.Glob("scel/*.scel")
	dictFile := "all.txt"
	var dictFileContent []string

	for _, scelFile := range scelFiles {
		records := getWordsFromSogouCellDict(scelFile)
		fmt.Printf("%s: %d 个词\n", scelFile, len(records))

		outFile := filepath.Join(dirs[0], strings.Replace(filepath.Base(scelFile), ".scel", ".txt", 1))
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			os.Exit(1)
		}
		defer f.Close()

		dictFileContent = append(dictFileContent, save(records, f)...)

		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Printf("合并后 %s: %d 个词\n", dictFile, len(dictFileContent))

	dictFileOut := filepath.Join(dirs[0], dictFile)
	fDict, err := os.Create(dictFileOut)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer fDict.Close()

	_, err = fDict.WriteString(strings.Join(dictFileContent, "\n"))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}
