package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"go-cleaner/checker"
	"go-cleaner/configurator"
)

const dumpFilePath = "./dump_file.txt"

var (
	dumpFile *os.File
	config   = &configurator.Config
)

func main() {
	configurator.Init()

	if !config.IsReady {
		log.Println("Not ready, check and adjust 'cleaner_config.yml' file to ger ready")
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

	log.Printf("Real mode %t \n", config.RealClean)

	log.Println("--->> Go...")
	err := checkAndRemove(config.StartPath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("<<--- Finished")
	printStats(config.RealClean)
	presentation()
}

func checkAndRemove(dirPath string) error {
	stats.FolderChecked++
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println(err)
		return err
	}

	printPath := true

	for _, entry := range entries {
		newPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() {
			checkAndRemove(newPath)
			// we don't need to check error here
			// if internal folder is absent it means that it's deleted or renamed
			// so we just skip this check
			continue
		}

		stats.FileChecked++

		if checker.IsNameMatch(entry.Name(), config.FileNames) {
			if config.RealClean {
				if config.RealClean {
					err := os.Remove(newPath)
					if err != nil {
						log.Println(err)
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
					log.Println(dirPath)
					printPath = false
				}
				log.Printf("   XXX: %s \n", entry.Name())
			}
		}

		if checker.IsExtMatch(entry, config.Exts) {
			if config.RealClean {
				err := os.Remove(newPath)
				if err != nil {
					log.Println(err)
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
				log.Println(dirPath)
				printPath = false
			}
			log.Printf("   XXX: %s \n", entry.Name())
		}
	}
	return nil
}
