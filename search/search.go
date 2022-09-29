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

// append with mutex
func (f *FileHits) Append(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Files = append(f.Files, path)
}

// TODO - parse xml properly
func SearchDocx(file string, str string) (bool, error) {

	f, err := zip.OpenReader(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	for _, file := range f.File {

		if strings.Compare(file.Name, "word/document.xml") == 0 {

			r, err := file.Open()
			if err != nil {
				return false, err
			}

			buf, err := io.ReadAll(r)
			if err != nil {
				return false, err
			}

			if strings.Contains(string(buf), str) {
				return true, nil
			}
		}
	}

	return false, nil
}

// return file path that is to be appended FileHits
func SearchTxt(file string, str string) (bool, error) {

	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	wordbuf, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(wordbuf), str) {
		return true, nil
	}

	return false, nil
}

func GrabStartDir(dir string) string {

	dir_size := len(dir)
	ret_buf := make([]byte, dir_size-1)

	if dir[dir_size-1] == '\\' || dir[dir_size-1] == '/' {

		copy(ret_buf, dir[:dir_size])

		return string(ret_buf)

	}

	return dir
}

func Run(Args []string) *FileHits {

	startingDir := GrabStartDir(Args[0])
	searchStr := Args[1]

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

				filePaths.Append(currPath)

			} else {

				fmt.Fprintln(os.Stderr, err)

			}
		} else if strings.Compare(filepath.Ext(currPath), ".txt") == 0 {

			if ok, err := SearchTxt(startingDir+"\\"+entry.Name(), searchStr); ok && err == nil {

				filePaths.Append(currPath)

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

				files.Append(path)

			}

		} else if strings.Compare(filepath.Ext(d.Name()), ".txt") == 0 && !d.IsDir() {

			if ok, err := SearchTxt(path, str); ok && err == nil {

				files.Append(path)

			}

		}

		return nil
	})

}
