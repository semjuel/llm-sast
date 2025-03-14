package android

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Checks if a file exists on disk.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Ensures a file is executable on Unix-like systems. (No-op on Windows).
func ensureExecutable(path string) error {
	// On Windows, .bat files don’t need +x permission; skip.
	if runtime.GOOS == "windows" {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	mode := info.Mode()
	// Check if “owner-executable” bit is already set; if not, set 0755
	if mode&0100 == 0 {
		return os.Chmod(path, 0755)
	}

	return nil
}

// Runs JADX with a given set of arguments. Redirects stdout/stderr to /dev/null.
// Returns the process exit code and any error encountered.
func runJadx(args []string) (int, error) {
	// Create a context with timeout from settings
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Send all stdout/stderr to /dev/null
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return -1, fmt.Errorf("failed to open /dev/null: %w", err)
	}
	defer devNull.Close()

	cmd.Stdout = devNull
	cmd.Stderr = devNull

	// Run the process
	err = cmd.Run()

	// If the context deadline was exceeded, interpret it as a timeout
	if ctx.Err() == context.DeadlineExceeded {
		return -1, fmt.Errorf("timeout running JADX: %w", ctx.Err())
	}

	// Retrieve exit code if possible
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), err
	}
	if err != nil {
		// Some other exec error (binary not found, permission denied, etc.)
		return -1, err
	}

	// Success, exit code 0
	if cmd.ProcessState != nil {
		return cmd.ProcessState.ExitCode(), nil
	}
	// Fallback if ProcessState is nil for some reason
	return 0, nil
}

// Apk2Java it attempts to decompile an APK with JADX.
// If that fails, it attempts to decompile all DEX files.
func Apk2Java(appPath, appDir string) error {
	// Build paths
	outputDir := filepath.Join(appDir, "java_source")

	// Remove existing output directory, ignoring errors
	_ = os.RemoveAll(outputDir)

	// Determine the correct JADX executable path
	// @TODO value is hardcoded
	jadxPath := "/usr/local/bin/jadx"

	// Ensure the file is executable (particularly on Unix)
	if err := ensureExecutable(jadxPath); err != nil {
		return err
	}

	// First attempt: decompile the APK
	args := []string{
		jadxPath, "-ds", outputDir,
		"-q", "-r", "--show-bad-code", appPath,
	}
	exitCode, err := runJadx(args)
	if err == nil && exitCode == 0 {
		// Success
		return nil
	}

	// If APK decompilation fails, attempt all *.dex files
	warnMsg := "Decompiling with JADX failed, attempting on all DEX files"
	log.Println(warnMsg)

	parentDir := filepath.Dir(appPath)
	decompileFailed := false

	// Walk all files under parentDir looking for *.dex
	filepath.Walk(parentDir, func(path string, info os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".dex") {
			msg := fmt.Sprintf("Decompiling %s with JADX", info.Name())
			log.Println(msg)

			// Modify the last argument to be the .dex file
			dexArgs := make([]string, len(args))
			copy(dexArgs, args)
			dexArgs[len(dexArgs)-1] = path

			dexExitCode, dexErr := runJadx(dexArgs)
			if dexErr != nil || dexExitCode != 0 {
				decompileFailed = true
				failMsg := fmt.Sprintf("Decompiling with JADX failed for %s", info.Name())
				log.Println(failMsg)
			}
		}
		return nil
	})

	if decompileFailed {
		failMsg := "Some DEX files failed to decompile"
		log.Println(failMsg)
	}

	// If we got this far, we either succeeded at least partially or failed
	if err != nil {
		// If the original APK attempt had a known error, return it
		return fmt.Errorf("apk2Java encountered an error: %w", err)
	}

	return nil
}
