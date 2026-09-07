package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	outputFlag := flag.String("o", "", "output archive name")
	flag.Parse()

	args := flag.Args()
	if *outputFlag == "" || len(args) == 0 {
		fmt.Println("Usage: r4go -o <archive.rar> <file_or_folder1> [file_or_folder2...]")
		os.Exit(1)
	}

	entries, err := CollectEntries(args)
	if err != nil {
		fmt.Printf("Error processing input paths: %v\n", err)
		os.Exit(1)
	}

	if err := CreateArchive(*outputFlag, entries); err != nil {
		fmt.Printf("Error creating archive: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated RAR: %s\n", *outputFlag)
}