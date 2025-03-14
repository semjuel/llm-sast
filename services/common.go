package services

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func unzip(r *zip.Reader, dest string) error {
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			// Create directories
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Create the necessary directories for the file
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Create the file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Open the file within the zip archive
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Copy the content to the output file
		_, err = io.Copy(outFile, rc)

		// Close the file within the zip archive
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
