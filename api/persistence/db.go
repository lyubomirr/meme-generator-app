package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	getOperationTimeout = time.Second * 2
	instance            *gorm.DB
	doOnce              sync.Once
)

func init() {
	err := migrate()
	if err != nil {
		panic(fmt.Sprintf("Couldn't migrate db: %v", err))
	}
	
	err = seedData()
	if err != nil {
		panic(fmt.Sprintf("couldn't seed data: %v", err))
	}
}

func migrate() error {
	db := getDB()
	err := db.AutoMigrate(&dbRole{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&dbUser{})
	if err != nil {
		return err
	}
	return nil
}

func getDB() *gorm.DB {
	if instance == nil {
		doOnce.Do(func() {
			dsn := "root:admin@tcp(127.0.0.1:3306)/memegenerator?charset=utf8mb4&parseTime=True&loc=Local"
			var err error
			instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				panic(err)
			}
		})
	}

	return instance
}

func seedData() error {
	db := getDB()

	result := db.FirstOrCreate(&dbRole{}, dbRole{
		ID: entities.AdminRoleId,
		Name:  entities.AdminRoleName,
	})
	if result.Error != nil {
		return result.Error
	}

	result = db.FirstOrCreate(&dbRole{}, &dbRole{
		ID: entities.NormalRoleId,
		Name: entities.NormalRoleName,
	})
	if result.Error != nil {
		return result.Error
	}

	var adminUser dbUser
	result = db.Where(&dbUser{
		Username: "admin",
		RoleID: entities.AdminRoleId,
	}).First(&adminUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = db.Create(&dbUser{
				Model:      gorm.Model{
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: gorm.DeletedAt{},
				},
				Username:   "admin",
				Password:   "$2a$10$4M/iO9uQzFmZW600Dcj39.Pv.K5E5IXR9zNwl0lpMKQajXwOPBp56",
				RoleID:     entities.AdminRoleId,
				PictureURL: "",
			})

			if result.Error != nil {
				return result.Error
			}
		} else {
			return result.Error
		}
	}
	return nil
}

func withReadTimeout(db *gorm.DB) *gorm.DB {
	return withTimeout(db, getOperationTimeout)
}

func withTimeout(db *gorm.DB, duration time.Duration) *gorm.DB {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return db.WithContext(ctx)
}

