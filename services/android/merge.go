package android

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/semjuel/llm-sast/utils"
	"os"
	"os/exec"
	"path/filepath"
)

func UnzipAPK(src, dest string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	ext := filepath.Ext(src)
	if ext == ".zip" {
		// Merge apk file inside the zip archive
		b, err = mergeApks(src)
		if err != nil {
			return err
		}
	}

	// Create a new zip.Reader from the in-memory []byte
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Iterate through each file/dir in the archive
	return utils.Unzip(r, dest)
}

func mergeApks(src string) ([]byte, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(b)
	zipReader, err := zip.NewReader(reader, int64(len(b)))
	if err != nil {
		return nil, err
	}

	dest := fmt.Sprintf("%s/%s", os.TempDir(), src)
	destZip := fmt.Sprintf("%s/zip", dest)
	_ = os.MkdirAll(destZip, 0755)
	defer os.RemoveAll(dest)

	err = utils.Unzip(zipReader, destZip)
	if err != nil {
		return nil, err
	}

	// Open (or create) a log file in append mode.
	logFile, err := os.OpenFile("/tmp/apkeditor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Merge split APKs
	mergedPath := filepath.Join(dest, "merged.apk")
	cmdApkEditor := exec.Command("apkeditor", "m", "-i", destZip, "-o", mergedPath)
	cmdApkEditor.Stdout = logFile
	cmdApkEditor.Stderr = logFile
	err = cmdApkEditor.Run()
	if err != nil {
		return nil, fmt.Errorf("apkeditor: %v", err)
	}

	signedPath := filepath.Join(dest, "signed-app.apk")
	keystorePath := "/keystore/upload-keystore.jks"
	// Validate if merged.apk exists
	if _, err := os.Stat(mergedPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("merged.apk not found: %s", mergedPath)
	}

	// Sign merged apk file
	cmdApkSigner := exec.Command("apksigner",
		"sign",
		"--ks", keystorePath,
		"--ks-key-alias", "upload",
		"--ks-pass", "pass:mypass",
		"--v3-signing-enabled", "true",
		"--out", signedPath,
		mergedPath,
	)
	cmdApkSigner.Stderr = os.Stderr
	err = cmdApkSigner.Run()
	if err != nil {
		return nil, fmt.Errorf("apksigner: %v", err)
	}

	return os.ReadFile(signedPath)
}
