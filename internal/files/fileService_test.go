package files

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	c "movie/system/internal/config" // Update with your actual module path
)

// mockConfig creates a dummy configuration
func mockConfig(basePath string) c.Config {
	return c.Config{
		Protocol: "http",
		HostName: "localhost",
		Port:     8080,
		File: c.ConfigFile{
			FilesBasePath: basePath,
			BucketName:    "test-bucket",
		},
	}
}
func TestGenerateUrl(t *testing.T) {
	tempDir := t.TempDir()
	fs := &FileService{Config: mockConfig(tempDir)}
	ctx := context.Background()

	bucketName := "avatars"
	objectKey := "a1b2-999999" // Shards will be "a1" and "b2"

	t.Run("success", func(t *testing.T) {
		// Setup the sharded directory correctly
		shard1, shard2 := objectKey[0:2], objectKey[2:4]
		dirPath := filepath.Join(tempDir, bucketName, shard1, shard2, objectKey)
		os.MkdirAll(dirPath, 0755)

		url, err := fs.GenerateUrl(ctx, bucketName, objectKey)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify the URL includes the shards
		expectedUrl := "http://localhost:8080/avatars/a1/b2/a1b2-999999"
		if url != expectedUrl {
			t.Errorf("expected %s, got %s", expectedUrl, url)
		}
	})

	t.Run("directory does not exist", func(t *testing.T) {
		_, err := fs.GenerateUrl(ctx, "avatars", "c3d4-00000")
		if err == nil {
			t.Fatal("expected error for non-existent directory, got nil")
		}
	})
}

func TestCreateFile(t *testing.T) {
	tempDir := t.TempDir()
	fs := &FileService{Config: mockConfig(tempDir)}
	ctx := context.Background()

	t.Run("successfully creates sharded file", func(t *testing.T) {
		bucketName := "uploads"
		fileName := "my doc.txt"
		fileContent := "hello world sharding"
		fileReader := strings.NewReader(fileContent)

		hasher := sha256.New()
		hasher.Write([]byte(fileContent))
		expectedHash := hex.EncodeToString(hasher.Sum(nil))

		result, err := fs.Create(ctx, bucketName, fileName, fileReader)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		if result.Hash != expectedHash {
			t.Errorf("expected hash %s, got %s", expectedHash, result.Hash)
		}

		// Verify URL format has 4 segments after domain: /bucket/shard1/shard2/objectKey
		if !strings.Contains(result.Url, "/uploads/"+result.ObjectKey[0:2]+"/"+result.ObjectKey[2:4]+"/") {
			t.Errorf("URL does not contain correct sharding structure: %s", result.Url)
		}

		// Verify file was actually saved in the sharded path
		savedFilePath := filepath.Join(
			tempDir,
			bucketName,
			result.ObjectKey[0:2],
			result.ObjectKey[2:4],
			result.ObjectKey,
			result.FileName,
		)

		data, err := os.ReadFile(savedFilePath)
		if err != nil {
			t.Fatalf("failed to read saved file at sharded path: %v", err)
		}

		if string(data) != fileContent {
			t.Errorf("expected file content %q, got %q", fileContent, string(data))
		}
	})

	t.Run("fails gracefully if base path is read-only", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping POSIX permission test on Windows")
		}

		readOnlyDir := filepath.Join(tempDir, "readonly")
		os.MkdirAll(readOnlyDir, 0400)

		badFs := &FileService{Config: mockConfig(readOnlyDir)}
		fileReader := strings.NewReader("some data")

		_, err := badFs.Create(ctx, "bucket", "file.txt", fileReader)
		if err == nil {
			t.Fatal("expected error when writing to read-only directory, got nil")
		}
	})
}
