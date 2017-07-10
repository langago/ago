package endpoint

import (
	"github.com/labstack/echo"
	"github.com/langago/ago/db"
	"net/http"
	"reflect"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/langago/ago/util"
)

type SimpleEndpoint struct {}

var Default Endpoint = nil

func (e SimpleEndpoint) Create(i interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(i).Elem()

	return func(c echo.Context) error {
		o := reflect.New(t).Interface()

		if err := c.Bind(o); err != nil {
			return err
		}

		if c.Echo().Validator != nil {
			if err := c.Validate(o); err != nil {
				return err
			}
		}

		if err := db.Debug().Create(o).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, o)
	}
}

func (e SimpleEndpoint) Get(i interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(i).Elem()

	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		v := reflect.New(t).Interface()

		if err := db.Debug().First(v, id).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, v)
	}
}

func (e SimpleEndpoint) List(i interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(i).Elem()

	queryFields := make(map[string]reflect.StructField)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		rangeTag := field.Tag.Get("query")
		if rangeTag != "" {
			queryFields[rangeTag] = field
		}
	}

	return func(c echo.Context) error {
		sliceType := reflect.SliceOf(t)
		//s := reflect.MakeSlice(sliceType, 0, 0).Interface()
		slice := reflect.New(sliceType).Interface()

		o := reflect.New(t).Interface()

		if err := c.Bind(o); err != nil {
			return err
		}

		d := db.Debug().Where(o)

		for field, structField := range queryFields {
			if p := c.QueryParam(field + "From"); p != "" {
				column := util.GetColumnName(structField)
				p, err := util.Apply(structField.Type.Kind(), p)
				if err != nil {
					return err
				}

				d = d.Where(column + " >= ?", p)
			}

			if p := c.QueryParam(field + "To"); p != "" {
				column := util.GetColumnName(structField)

				p, err := util.Apply(structField.Type.Kind(), p)
				if err != nil {
					return err
				}

				d = d.Where(column + " < ?", p)
			}
		}

		if err := d.Find(slice).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, slice)
	}
}

func (e SimpleEndpoint) Update(i interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(i).Elem()

	var omitFields []string
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("update") == "-" {
			column := t.Field(i).Tag.Get("column")
			if column == "" {
				column = gorm.ToDBName(t.Field(i).Name)
			}
			omitFields = append(omitFields, column)
		}
	}

	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		v := reflect.New(t)
		v.Elem().FieldByName("ID").SetUint(uint64(id))

		o := reflect.New(t).Interface()

		if err := c.Bind(o); err != nil {
			return err
		}

		if err := db.Debug().Model(v.Interface()).Omit(omitFields...).Updates(o).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, o)
	}
}

func (e SimpleEndpoint) Delete(i interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(i).Elem()

	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		v := reflect.New(t)
		v.Elem().FieldByName("ID").SetUint(uint64(id))

		if err := db.Debug().Delete(v.Interface()).Error; err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}