package service

import (
	"io"
	"os"
	"path/filepath"
)

// SaveFile menyimpan file ke local storage
func SaveFile(basePath, nis, fileName string, reader io.Reader) (string, error) {

	// folder siswa
	dir := filepath.Join(basePath, nis)

	// buat folder jika belum ada
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// full path file
	fullPath := filepath.Join(dir, fileName)

	// buat file tujuan
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy isi file upload ke file tujuan
	if _, err := io.Copy(dst, reader); err != nil {
		return "", err
	}

	// return relative path untuk disimpan ke DB
	relativePath := filepath.Join(nis, fileName)

	return relativePath, nil
}

// DeleteFile menghapus file lama
func DeleteFile(basePath, relativePath string) error {

	if relativePath == "" {
		return nil
	}

	fullPath := filepath.Join(basePath, relativePath)

	// hapus file
	if err := os.Remove(fullPath); err != nil {

		// kalau file tidak ada, ignore
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return nil
}
