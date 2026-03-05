package files

import (
	"context"
	"io"
)

type IFiles interface {
	//Create new file and store it according to your choice
	//Saving the created file to the database at files table with failure handling for the saved file
	//for filesService bucketName is the main folder and objectKey is the secondary folder
	//A random objectKey will be chosen for the file and the original file name will be saved to the db
	Create(ctx context.Context, bucketName string, fileName string, file io.Reader) (*File, error)

	//Delete the file if exists if not it will return error as nil
	//returns error if the db delete process failed
	Delete(ctx context.Context, bucketName string, objectKey string) error

	//NOT IMPLEMENTED
	GetByKeys(ctx context.Context, bucketName string, objectKey []string) ([]File, error)

	//Generate URL for the file to be accessed for the internet
	//Implementation is specific to the service used
	//FileService: API link to the app (NOT IMPLEMENTED)
	//S3Service: A presigned URL valid for 1 hour
	GenerateUrl(ctx context.Context, bucketName string, objectKey string) (string, error)
}
