package main

import "log"

const version = "2.2.0"

var stats Stats

type Stats struct {
	FolderChecked uint16
	FileChecked   uint16
	FoundCount    uint16
	RemovedCount  uint16
}

func presentation() {
	text := "Version: %s * Developped by Oleh Sobchuk tel: 0730240643\n"
	log.Printf(text, version)
}

func printStats(realClean bool) {
	log.Printf("\nChecked: %d folder(s), %d file(s)\n", stats.FileChecked, stats.FileChecked)
	if realClean {
		log.Printf("\nRemoved: %d file(s)\n\n", stats.RemovedCount)
	} else {
		log.Printf("\nFound: %d file(s)\n\n", stats.FoundCount)
	}
}
