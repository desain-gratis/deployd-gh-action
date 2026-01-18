package utility

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func IsDir(input string) (bool, error) {
	info, err := os.Stat(input)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, nil
	}

	return true, nil
}

func BundleDir(output, srcDir string) error {
	if err := ensureParentDir(output); err != nil {
		panic(err)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	base := filepath.Base(srcDir)

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil || rel == "." {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		hdr.Name = filepath.ToSlash(filepath.Join(base, rel))

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(tw, f)
			return err
		}

		return nil
	})
}

type Result struct {
	Path string
	Size int64
	Hash string
}

func Calculate(path string) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	size, err := io.Copy(h, f)
	if err != nil {
		return nil, err
	}

	return &Result{
		Path: path,
		Size: size,
		Hash: hex.EncodeToString(h.Sum(nil)),
	}, nil
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
