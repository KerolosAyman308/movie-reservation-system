package files

type File struct {
	ObjectKey    string `gorm:"primaryKey;type: varchar(40)"`
	BucketName   string `gorm:"not null;type: varchar(250)"`
	OriginalName string `gorm:"not null;type: varchar(250)"`
	FileName     string `gorm:"not null;type: varchar(250)"`
	Url          string `gorm:"-"`
	Size         int64  `gorm:"not null;"`
	Hash         string `gorm:"not null;type: varchar(65)"`
}
