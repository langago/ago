package db

import "github.com/jinzhu/gorm"

var DB *gorm.DB

func AutoMigrate(value ...interface{}) *gorm.DB {
	return DB.AutoMigrate(value...)
}

func Debug() *gorm.DB {
	return DB.Debug()
}

func Create(value interface{}) *gorm.DB {
	return DB.Create(value)
}

func First(out interface{}, where ...interface{}) *gorm.DB {
	return DB.First(out, where...)
}

func Where(query interface{}, args ...interface{}) *gorm.DB {
	return DB.Where(query, args...)
}

func Find(out interface{}, where ...interface{}) *gorm.DB {
	return DB.Find(out, where...)
}