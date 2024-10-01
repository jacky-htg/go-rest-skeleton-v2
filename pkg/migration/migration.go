package migration

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type walk filepath.WalkFunc

func Migrate(db *sql.DB) error {
	// Buat tabel migrations jika belum ada
	err := createMigrationsTable(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	var mywalk walk = func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".sql" {
			err := executeSQLFile(tx, path)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return nil
	}

	err = filepath.Walk("migration", filepath.WalkFunc(mywalk))
	if err != nil {
		return fmt.Errorf("error walking the path: %v", err)
	}

	return tx.Commit()
}

func createMigrationsTable(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS _migrations (
        id SERIAL PRIMARY KEY,
        filename TEXT NOT NULL,
        checksum TEXT NOT NULL,
        executed_at TIMESTAMPTZ DEFAULT timezone('utc', now())
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("could not create migrations table: %v", err)
	}

	return nil
}

// executeSQLFile membaca dan mengeksekusi file SQL
func executeSQLFile(tx *sql.Tx, path string) error {
	// Hitung checksum file SQL
	checksum, err := calculateChecksum(path)
	if err != nil {
		return fmt.Errorf("could not calculate checksum for file %s: %v", path, err)
	}

	// Periksa apakah file sudah dieksekusi sebelumnya
	var exists bool
	var storedChecksum string
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM _migrations WHERE filename = $1), checksum 
		FROM _migrations 
		WHERE filename = $1
	`, filepath.Base(path)).Scan(&exists, &storedChecksum)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("could not check migrations: %v", err)
	}

	if exists {
		if storedChecksum != checksum {
			return fmt.Errorf("checksum mismatch for file %s; file has changed", filepath.Base(path))
		}
		return nil
	}

	// Baca isi file SQL dan eksekusi
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", filepath.Base(path), err)
	}

	_, err = tx.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute SQL file %s: %v", filepath.Base(path), err)
	}

	// Insert ke tabel migrations
	_, err = tx.Exec("INSERT INTO _migrations (filename, checksum) VALUES ($1, $2)", filepath.Base(path), checksum)
	if err != nil {
		return fmt.Errorf("could not insert migration record for file %s: %v", filepath.Base(path), err)
	}

	return nil
}

// calculateChecksum menghitung checksum SHA256 dari file
func calculateChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
