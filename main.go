package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	instMCDir := os.Getenv("INST_MC_DIR")
	var sourceDir string
	destDir := "./screenshots"

	if instMCDir != "" {
		sourceDir = filepath.Join(instMCDir, "screenshots")
	}

	flag.StringVar(&sourceDir, "from", sourceDir, "Instance screenshots directory")
	flag.StringVar(&destDir, "to", destDir, "Destination directory for the screenshots")
	flag.Parse()

	if sourceDir == "" {
		log.Fatal("Instance screenshots directory is not specified.")
	}

	sameVolume := filepath.VolumeName(sourceDir) == filepath.VolumeName(destDir)

	filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			err := fmt.Errorf("failed to access %v", path)
			log.Printf("An error occured while accesing source directory: %v", err)
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		dest := filepath.Join(destDir, d.Name())

		if sameVolume {
			err = os.Rename(path, dest)
		} else {
			err = move(path, dest)
		}

		if err != nil {
			err := fmt.Errorf("failed to move %v to %v", path, dest)
			log.Printf("An error occured while moving files: %v", err)
			return nil
		}

		log.Printf("Moved %v to %v.\n", path, dest)
		return nil
	})
}

func move(src string, dest string) error {
	inputFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("failed to copy: %v", err)
	}

	inputFile.Close()

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %v", err)
	}
	return nil
}
