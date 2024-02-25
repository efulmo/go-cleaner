package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const dumpFilePath = "./dump_file.txt"

var (
	dumpFile *os.File
)

func main() {
	SetConfiguration()

	if !config.IsReady {
		fmt.Println("Not ready, check config file")
		return
	}

	if !config.RealClean {
		var err error
		dumpFile, err = os.Create(dumpFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer dumpFile.Close()
		defer dumpFile.Sync()
	}

	fmt.Printf("Real mode %t \n", config.RealClean)

	fmt.Println("--->> Go...")
	checkAndRemove(config.StartPath)
	fmt.Println("<<--- Finished")
	printStats()
	presentation()
}

func checkAndRemove(dirPath string) {
	stats.FolderChecked++
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	printPath := true

	for _, entry := range entries {
		newPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() {
			checkAndRemove(newPath)
			continue
		}

		stats.FileChecked++

		if isSuspicious(entry) {
			if config.RealClean {
				err := os.Remove(newPath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				stats.RemovedCount++
			} else {
				stats.FoundCount++
				_, err = dumpFile.WriteString(fmt.Sprintf("%s\n", newPath))
				if err != nil {
					log.Println("Write to file error:", err)
				}
			}
			if printPath {
				fmt.Println(dirPath)
				printPath = false
			}
			fmt.Printf("   XXX: %s \n", entry.Name())
		}
	}
}

func isSuspicious(file fs.DirEntry) bool {
	ext := extension(file.Name())
	return slices.Contains(config.Exts, ext)
}

func extension(fileName string) string {
	ext := filepath.Ext(fileName)
	ext = strings.TrimPrefix(ext, ".")
	return strings.ToLower(ext)
}
