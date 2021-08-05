package data

import (
	"os"

	"github.com/go-playground/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Post struct {
		gorm.Model
		PostBody     string  `gorm:"not null" json:"post_body" validate:"required"`
		Author       string  `gorm:"not null" json:"author" validate:"required"`
		ImageURI     string  `json:"image_uri"`
		Likes        int     `gorm:"default:0" json:"num_likes"`
		RefersToPost *uint   `json:"refers_to_post"`
		Comments     []*Post `gorm:"foreignkey:RefersToPost" json:"comments"`
	}

	PostRepo struct {
		db *gorm.DB
	}
)

func NewPostRepo() *PostRepo {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Post{})
	return &PostRepo{db}
}

func (pr *PostRepo) CreatePost(post *Post) error {
	return pr.db.Create(post).Error
}

func (p *Post) Validate() error {
	validator := validator.New()
	return validator.Struct(p)
}
