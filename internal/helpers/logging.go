package helpers

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"syscall"
	"time"
)

func LogDir() string {
	return GetEnvFallback("LOG_DIR", "logs")
}

func SetupLogging() (*os.File, error) {
	dir := LogDir()

	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	path := fmt.Sprintf("app-%d.log", time.Now().Unix())

	f, err := root.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}

	console, err := dup(os.Stdout)
	if err != nil {
		console = os.Stdout
	}

	writer := io.MultiWriter(console, f)
	log.SetOutput(writer)
	slog.SetDefault(slog.New(slog.NewTextHandler(writer, nil)))

	pipeStdStreams(console, f)

	return f, nil
}

func dup(f *os.File) (*os.File, error) {
	fd, err := syscall.Dup(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	return os.NewFile(uintptr(fd), f.Name()), nil
}

func pipeStdStreams(console, file *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		return
	}

	if err := syscall.Dup2(int(w.Fd()), int(os.Stdout.Fd())); err != nil {
		return
	}
	if err := syscall.Dup2(int(w.Fd()), int(os.Stderr.Fd())); err != nil {
		return
	}

	os.Stdout, os.Stderr = w, w

	go func() {
		_, _ = io.Copy(io.MultiWriter(console, file), r)
	}()
}
