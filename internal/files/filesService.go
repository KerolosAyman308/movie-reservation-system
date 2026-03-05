package files

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	c "movie/system/internal/config"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type FileService struct {
	Config c.Config
	DB     *gorm.DB
}

func NewFileService(config c.Config, db *gorm.DB) IFiles {
	return &FileService{Config: config, DB: db}
}

// Create implements IFiles.
func (f *FileService) Create(ctx context.Context, bucketName string, fileName string, file io.Reader) (*File, error) {
	// Generate object key
	objectKey := fmt.Sprintf("%s", strings.ToLower(rand.Text()))
	shard1 := objectKey[0:2]
	shard2 := objectKey[2:4]
	baseUrl := filepath.Join(f.Config.File.FilesBasePath, bucketName, shard1, shard2, objectKey)
	//Create the object key so the the structure will be /bucketName/objectKey/folderName
	err := os.MkdirAll(baseUrl, 0755)
	if err != nil {
		return nil, err
	}

	// refine the file to be compatible
	fileNameToSave := RefineS3Filename(fileName)
	path := filepath.Join(f.Config.File.FilesBasePath, bucketName, shard1, shard2, objectKey, fileNameToSave)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	//hash the file to check if it exists later
	hasher := sha256.New()
	//multi write to the hash and the file
	multiWriter := io.MultiWriter(out, hasher)
	size, err := io.Copy(multiWriter, file)
	if err != nil {
		os.RemoveAll(baseUrl)
		return nil, err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	url, _ := f.GenerateUrl(ctx, bucketName, objectKey)
	fileToSave := &File{
		OriginalName: fileName,
		FileName:     fileNameToSave,
		ObjectKey:    objectKey,
		Size:         size,
		Hash:         hash,
		Url:          url,
		BucketName:   bucketName,
	}

	err = f.DB.Model(&File{}).WithContext(ctx).Create(fileToSave).Error
	if err != nil {
		os.RemoveAll(baseUrl)
		return nil, err
	}
	return fileToSave, nil
}

// Delete implements IFiles.
func (f *FileService) Delete(ctx context.Context, bucketName string, objectKey string) error {
	shard1 := objectKey[0:2]
	shard2 := objectKey[2:4]

	baseUrl := filepath.Join(f.Config.File.FilesBasePath, bucketName, shard1, shard2, objectKey)
	err := f.DB.Model(&File{}).WithContext(ctx).Delete(&File{}, "object_key = ?", objectKey).Error
	if err != nil {
		return err
	}
	return os.RemoveAll(baseUrl)
}

// GenerateUrl implements IFiles.
func (f *FileService) GenerateUrl(ctx context.Context, bucketName string, objectKey string) (string, error) {
	shard1 := objectKey[0:2]
	shard2 := objectKey[2:4]
	baseUrl := filepath.Join(f.Config.File.FilesBasePath, bucketName, shard1, shard2, objectKey)
	info, err := os.Stat(baseUrl)
	if errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if !info.IsDir() {
		return "", ErrExpectedDir
	}

	filePath := fmt.Sprintf("%s/%s/%s/%s", bucketName, shard1, shard2, objectKey)
	return fmt.Sprintf("%s://%s:%d/%s", f.Config.Protocol, f.Config.HostName, f.Config.Port, filePath), nil
}

// GetByKeys implements IFiles.
func (f *FileService) GetByKeys(ctx context.Context, bucketName string, objectKey []string) ([]File, error) {
	panic("unimplemented")
}
