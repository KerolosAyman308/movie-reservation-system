package config

type ConfigFile struct {
	BucketName    string
	FilesBasePath string
	AWSAccessKey  string
	AWSSecretKey  string
	AWSHost       string
	UseFile       bool
}
