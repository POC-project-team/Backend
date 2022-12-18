package postgres

import (
	"backend/internal/controller/rest/response"
	"backend/internal/entity"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Client struct {
	db *gorm.DB
}

func NewClient() (*Client, error) {
	client := &Client{}
	db, err := client.Connect()
	if err != nil {
		return nil, err
	}
	client.db = db
	return client, nil
}

func (c *Client) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=admin password=admin dbname=poc_db port=5432"))
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&entity.Tag{}); err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) GetUser(login, password string) (entity.User, error) {
	var user entity.User
	if err := c.db.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (c *Client) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	if err := c.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) CreateUser(login, password string) (entity.User, error) {
	var user entity.User

	// check if user with such login already exists
	if err := c.db.Where("login = ?", login).First(&user).Error; err == nil {
		return user, errors.New(response.UserAlreadyExists)
	}

	// create new user with unique id
	if err := c.db.Create(&entity.User{Login: login, Password: password}).Error; err != nil {
		return entity.User{
			UserID: 0,
		}, err
	}

	// add and return new user
	if err := c.db.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return entity.User{
			UserID: 0,
		}, err
	}

	return entity.User{UserID: user.UserID}, nil
}

func (c *Client) ChangePassword(userId uint, password string) error {
	if err := c.db.Model(&entity.User{}).Where("user_id = ?", userId).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ChangeLogin(userId uint, login string) error {
	if err := c.db.Model(&entity.User{}).Where("user_id = ?", userId).Update("login", login).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteUser(userId uint) error {
	if err := c.db.Where("user_id = ?", userId).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) GetUserTags(userId uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	if err := c.db.Where("user_id = ?", userId).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *Client) CreateTag(userId uint, tagId, tagName string) (entity.Tag, error) {
	var tag entity.Tag
	if err := c.db.Create(&entity.Tag{UserID: userId, TagID: tagId, TagName: tagName}).Error; err != nil {
		return entity.Tag{}, err
	}
	// return the created tag
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).First(&tag).Error; err != nil {
		return entity.Tag{}, err
	}
	return tag, nil
}

func (c *Client) GetTag(userId uint, tagId string) (entity.Tag, error) {
	var tag entity.Tag

	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).First(&tag).Error; err != nil {
		return entity.Tag{}, errors.New(response.NoSuchTag)
	}
	return tag, nil
}

func (c *Client) DeleteTag(userId uint, tagId string) error {
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).Delete(&entity.Tag{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTag(userId uint, tagId, tagName string) (entity.Tag, error) {
	var tag entity.Tag
	if err := c.db.Model(&entity.Tag{}).Where("user_id = ? AND tag_id = ?", userId, tagId).Update("tag_name", tagName).Error; err != nil {
		return tag, errors.New(response.NoSuchTag)
	}
	return tag, nil
}

func (c *Client) GetUserNotes(userId uint, tagId string) ([]entity.Note, error) {
	// check the user
	var user entity.User
	if err := c.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New(response.UserNotFound)
	}
	var notes []entity.Note
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (c *Client) TransferTag(userId uint, tagId, login string) error {
	var user entity.User
	if err := c.db.Where("login = ?", login).First(&user).Error; err != nil {
		return errors.New(response.UserNotFound)
	}
	if err := c.db.Model(&entity.Tag{}).Where("user_id = ? AND tag_id = ?", userId, tagId).Update("user_id", user.UserID).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) AddNote(userId uint, tagId, noteInfo string) (entity.Tag, error) {
	var tag entity.Tag

	if err := c.db.Create(&entity.Note{UserId: userId, TagID: tagId, Note: noteInfo, Time: time.Now()}).Error; err != nil {
		return tag, errors.New(response.NoSuchTag)
	}
	return tag, nil
}
