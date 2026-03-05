package files

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	log "log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

type S3Service struct {
	S3Client      *s3.Client
	DB            *gorm.DB
	S3Manager     *manager.Uploader
	PresignClient *s3.PresignClient
}

func NewS3Service(s *s3.Client, db *gorm.DB, manager *manager.Uploader, signer *s3.PresignClient) IFiles {
	return &S3Service{
		S3Client:      s,
		DB:            db,
		S3Manager:     manager,
		PresignClient: signer,
	}
}

// Create implements IFiles.
func (s *S3Service) Create(ctx context.Context, bucketName string, fileName string, file io.Reader) (*File, error) {
	objectKey := fmt.Sprintf("%s", strings.ToLower(rand.Text()))
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = s3.NewObjectExistsWaiter(s.S3Client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		log.Error(err.Error(), "Failed attempt to wait for object to exist.\n", objectKey)
		return nil, err
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to rewind file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("file does not support seeking")
	}

	hash := sha256.New()
	sizeBytes, err := io.Copy(hash, file) // io.Copy returns the total bytes read!
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	url, _ := s.GenerateUrl(ctx, bucketName, objectKey)
	fileToSave := &File{
		OriginalName: fileName,
		FileName:     fileName,
		ObjectKey:    objectKey,
		Size:         sizeBytes,
		Hash:         hashString,
		Url:          url,
		BucketName:   bucketName,
	}

	err = s.DB.Model(&File{}).WithContext(ctx).Create(fileToSave).Error
	if err != nil {
		s.Delete(ctx, bucketName, objectKey)
		return nil, err
	}
	return fileToSave, nil
}

// Delete implements IFiles.
func (s *S3Service) Delete(ctx context.Context, bucketName string, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	input.VersionId = aws.String("")
	input.BypassGovernanceRetention = aws.Bool(true)

	err := s.DB.Model(&File{}).WithContext(ctx).Delete(&File{}, "object_key = ?", objectKey).Error
	if err != nil {
		return err
	}

	_, err = s.S3Client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}
	err = s3.NewObjectNotExistsWaiter(s.S3Client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		return err
	}
	return nil
}

// GenerateUrl implements IFiles.
func (s *S3Service) GenerateUrl(ctx context.Context, bucketName string, objectKey string) (string, error) {
	request, err := s.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * 60 * int64(time.Second))
	})

	if err != nil {
		return "", err
	}
	return request.URL, err
}

// GetByKeys implements IFiles.
func (s *S3Service) GetByKeys(ctx context.Context, bucketName string, objectKey []string) ([]File, error) {
	panic("unimplemented")
}
