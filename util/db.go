package util

import (
	"reflect"
	"github.com/jinzhu/gorm"
)

func GetColumnName(s reflect.StructField) string {
	c := s.Tag.Get("column")

	if c == "" {
		c = gorm.ToDBName(s.Name)
	}

	return c
}