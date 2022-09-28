package search

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileHits struct {
	Files []string
	mu    sync.Mutex
}

func NewFileHits() *FileHits {
	return new(FileHits)
}

// TODO
func SearchDocx(file string, str string) (bool, error) {
	r, err := zip.OpenReader(file)
	if err != nil {
		return false, err
	}
	defer r.Close()
	return false, nil
}

// return file path that is to be appended FileHits
func SearchTxt(file string, str string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	wordbuf, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(wordbuf), str) {
		return true, nil
	}

	return false, nil
}

func Run(Args []string) *FileHits {

	startingDir := Args[1]
	searchStr := Args[2]

	// need to do checking to see if startingDir ends with a '\'
	dirEntries, err := os.ReadDir(startingDir + "\\")
	if err != nil {
		log.Fatal(err)
	}

	filePaths := NewFileHits()

	var wg sync.WaitGroup

	// loop through the root directory and create a thread for every directory
	// append a .txt or .docx file to the file list if said document contains
	// string user is searching for
	for _, entry := range dirEntries {

		currPath := startingDir + "\\" + entry.Name()

		if entry.IsDir() {

			wg.Add(1)
			go SearchDir(currPath, searchStr, filePaths, &wg)

		} else if strings.Compare(filepath.Ext(currPath), ".docx") == 0 {

			if ok, err := SearchDocx(startingDir+"\\"+entry.Name(), searchStr); ok && err == nil {

				filePaths.mu.Lock()
				filePaths.Files = append(filePaths.Files, currPath)
				filePaths.mu.Unlock()

			} else {

				fmt.Fprintln(os.Stderr, err)

			}
		} else if strings.Compare(filepath.Ext(currPath), ".txt") == 0 {

			if ok, err := SearchTxt(startingDir+"\\"+entry.Name(), searchStr); ok && err == nil {

				filePaths.mu.Lock()
				filePaths.Files = append(filePaths.Files, currPath)
				filePaths.mu.Unlock()

			} else {

				fmt.Fprintln(os.Stderr, err)

			}
		}
	}

	wg.Wait()

	return filePaths
}

// looks at folder, file names, and strings within .docx
func SearchDir(dir string, str string, files *FileHits, wg *sync.WaitGroup) {

	defer wg.Done()

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

		if strings.Compare(filepath.Ext(d.Name()), ".docx") == 0 && !d.IsDir() {

			if ok, err := SearchDocx(path, str); ok && err == nil {

				files.mu.Lock()
				files.Files = append(files.Files, path)
				files.mu.Unlock()

			}

		} else if strings.Compare(filepath.Ext(d.Name()), ".txt") == 0 && !d.IsDir() {

			if ok, err := SearchTxt(path, str); ok && err == nil {

				files.mu.Lock()
				files.Files = append(files.Files, path)
				files.mu.Unlock()

			}

		}

		return nil
	})

}
