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
		Author       string  `gorm:"not null" json:"author"`
		ImageURI     string  `json:"image_uri"`
		Likes        int     `gorm:"default:0" json:"num_likes"`
		Reported     bool    `gorm:"default:0" json:"is_reported"`
		RefersToPost *uint   `json:"refers_to_post"`
		Comments     []*Post `gorm:"foreignkey:refers_to_post" json:"comments"`
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

func (pr *PostRepo) CreatePost(post *Post, id uint) error {
	if id != 0 {
		p := pr.GetPost(id)
		if p.ID != 0 {
			post.RefersToPost = &id
			return pr.db.Create(post).Error
		}
	}
	return pr.db.Create(post).Error
}

func (pr *PostRepo) GetPost(id uint) *Post {
	p := Post{}
	pr.db.Preload("Comments").Where("id = ?", id).First(&p)
	return &p
}

func (p *Post) Validate() error {
	validator := validator.New()
	return validator.Struct(p)
}
