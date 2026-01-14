package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.Printf("Started prismahelper")

	var sourceDir, destDir string

	flag.StringVar(&sourceDir, "from", sourceDir, "Instance screenshots directory")
	flag.StringVar(&destDir, "to", destDir, "Destination directory for the screenshots")
	flag.Parse()

	if sourceDir == "" {
		log.Printf("Source directory not specified, getting default value...")

		instMCDir := os.Getenv("INST_MC_DIR")
		if instMCDir == "" {
			log.Fatalf("Failed to get default source directory: environment variable INST_MC_DIR has no value")
		}

		sourceDir = filepath.Join(instMCDir, "screenshots")
	}

	if destDir == "" {
		log.Printf("Destination directory not specified, getting default value...")

		execFile, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get default destination directory: failed to get path of executable file: %v", err)
		}

		execDir := filepath.Dir(execFile)
		destDir = filepath.Join(execDir, "screenshots")
	}

	sourceAbs, err := filepath.Abs(sourceDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of source directory %v: %v", sourceDir, err)
	}

	destAbs, err := filepath.Abs(destDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of destination directory %v: %v", destDir, err)
	}

	log.Printf("Source directory: %v", sourceAbs)
	log.Printf("Destination directory: %v", destAbs)

	err = os.MkdirAll(destAbs, 0755)
	if err != nil {
		log.Fatalf("Failed to make directories for the destination %v: %v", destAbs, err)
	}

	filepath.WalkDir(sourceAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Failed to access %v: %v", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		dest := filepath.Join(destAbs, d.Name())
		dest = uniquePath(dest)

		err = safeMove(path, dest)
		if err != nil {
			log.Fatalf("Failed to move %v -> %v: %v", path, dest, err)
			return nil
		}

		log.Printf("Moved %v -> %v", path, dest)
		return nil
	})
}

func safeMove(source, dest string) error {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of source directory %v: %v", source, err)
	}

	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of destination directory %v: %v", dest, err)
	}

	if filepath.VolumeName(sourceAbs) == filepath.VolumeName(destAbs) {
		err = os.Rename(sourceAbs, destAbs)
		if err == nil {
			return nil
		}
	}

	sourceInfo, err := os.Stat(sourceAbs)
	if err != nil {
		return fmt.Errorf("failed to get file information of source file %v: %v", sourceAbs, err)
	}
	expectedSize := sourceInfo.Size()

	inputFile, err := os.Open(sourceAbs)
	if err != nil {
		return fmt.Errorf("failed to open source file %v: %v", sourceAbs, err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destAbs)
	if err != nil {
		return fmt.Errorf("failed to create destination file %v: %v", destAbs, err)
	}

	bytesCopied, err := io.Copy(outputFile, inputFile)
	if err != nil {
		outputFile.Close()
		os.Remove(destAbs)
		return fmt.Errorf("failed to copy %v -> %v: %v", inputFile, outputFile, err)
	}

	err = outputFile.Sync()
	if err != nil {
		outputFile.Close()
		os.Remove(destAbs)
		return fmt.Errorf("failed to sync destination file %v with the filesystem: %v", outputFile, err)
	}

	err = outputFile.Close()
	if err != nil {
		os.Remove(destAbs)
		return fmt.Errorf("failed to close destination file %v: %v", outputFile, err)
	}

	if bytesCopied != expectedSize {
		os.Remove(destAbs)
		return fmt.Errorf("copy verification failed for file %v (expected %d bytes, got %d)", 
			outputFile, expectedSize, bytesCopied)
	}

	return os.Remove(sourceAbs)
}

func uniquePath(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return path
	}

	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s_%d%s", name, i, ext)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(newPath); errors.Is(err, fs.ErrNotExist) {
			return newPath
		}
	}
}
