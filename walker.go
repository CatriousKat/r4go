package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileEntry struct {
	DiskPath    string
	ArchiveName string
}

func CollectEntries(inputPaths []string) ([]FileEntry, error) {
	var entries []FileEntry

	for _, arg := range inputPaths {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to access %s: %v", arg, err)
		}

		if !info.IsDir() {
			entries = append(entries, FileEntry{
				DiskPath:    arg,
				ArchiveName: filepath.ToSlash(arg),
			})
		} else {
			parentDir := filepath.Dir(arg)
			err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(parentDir, path)
				if err != nil {
					rel = path
				}
				entries = append(entries, FileEntry{
					DiskPath:    path,
					ArchiveName: filepath.ToSlash(rel),
				})
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed walking directory %s: %v", arg, err)
			}
		}
	}

	return entries, nil
}