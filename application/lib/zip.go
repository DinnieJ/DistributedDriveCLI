package lib

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func AddZip(src, dest string) error {
	var zipFile, err = os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	var zipWriter = zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println("Zipping up:", path)
		if info.IsDir() {
			path = fmt.Sprintf("%s%c", path, os.PathSeparator)
			_, err = zipWriter.Create(path)
			return err
		}
		var file, e = os.Open(path)
		if e != nil {
			return err
		}
		defer file.Close()

		var f, e2 = zipWriter.Create(path)
		if e2 != nil {
			return err
		}

		var _, errCopy = io.Copy(f, file)
		if errCopy != nil {
			return errCopy
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
