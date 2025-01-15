package main

import (
	"log"
	"os"
	"flag"
	"path/filepath"
	"io/fs"
)

func main() {
	inst_mc_dir := os.Getenv("INST_MC_DIR")
	var source_dir string
	dest_dir := "./screenshots"

	if inst_mc_dir != "" {
		source_dir = filepath.Join(inst_mc_dir, "screenshots") 
	}
	
	flag.StringVar(&source_dir, "from", source_dir, "Instance screenshots directory")
	flag.StringVar(&dest_dir, "to", dest_dir, "Destination directory for the screenshots")
	flag.Parse()

	if source_dir == "" {
		log.Fatal("Instance screenshots directory is not specified.")
	}

	filepath.WalkDir(source_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Encountered an error while accesing %v.\n", path)
			return filepath.SkipDir
		}
	
		if d.IsDir() {
			return nil
		}

		dest := filepath.Join(dest_dir, d.Name())

		err = os.Rename(path, dest)
		if err != nil {
			log.Printf("Encountered an error while moving %v to %v.\n", path, dest)
			return nil
		}

		log.Printf("Moved %v to %v.\n", path, dest)
		return nil
	})
}
