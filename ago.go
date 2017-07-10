package ago

import (
	"github.com/labstack/echo"
	"reflect"
	strs "github.com/langago/going/strings"
	"github.com/jinzhu/inflection"
	"github.com/langago/ago/db"
	"github.com/langago/ago/endpoint"
	"strings"
	log "github.com/sirupsen/logrus"
)


type Ago struct {
	*echo.Echo
}

type ComponentRegistry struct {
	ModelType reflect.Type
	endpoint.Endpoint
}

const (
	ENTITY_CREATE = "Create"
	ENTITY_GET = "Get"
	ENTITY_LIST = "List"
	ENTITY_UPDATE = "Update"
	ENTITY_DELETE = "Delete"
)

var registries map[string]ComponentRegistry

func New() (a *Ago) {
	a = &Ago{
		echo.New(),
	}

	a.HideBanner = true

	registries = make(map[string]ComponentRegistry)

	return
}

func (a *Ago) Register(i interface{}, e endpoint.Endpoint) {
	// create table or update fields
	db.AutoMigrate(i)

	t := reflect.TypeOf(i).Elem()
	typeName := t.Name()
	tnVar := strs.LowerFirstChar(typeName)

	// registry := ComponentRegistry{t,endpoint.SimpleEndpoint{}}
	var registry ComponentRegistry
	if e != nil {
		registry = ComponentRegistry{t,e}
	} else {
		registry = ComponentRegistry{t,endpoint.SimpleEndpoint{}}
	}
	registries[tnVar] = registry

	a.registerHandlers(tnVar)
}

func (a *Ago) registerHandlers(typeName string) {
	// pluralizes
	groupUri := inflection.Plural(typeName)
	group := a.Group("/" + groupUri)

	registry := registries[typeName]

	v := reflect.ValueOf(registry.Endpoint)
	for i := 0; i < v.NumMethod(); i++ {
		methodName := v.Type().Method(i).Name
		//v1 := v.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(typeName), reflect.ValueOf(registry.Model)})
		v1 := v.MethodByName(methodName).Call([]reflect.Value{reflect.New(registry.ModelType)})
		fn := v1[0].Interface()

		if strings.HasPrefix(methodName, ENTITY_CREATE) {
			subUri := strings.ToLower(strings.TrimPrefix(methodName, ENTITY_CREATE))
			group.POST(subUri, fn.(echo.HandlerFunc))
			log.Debugf("Mapped POST /%s", groupUri)
		} else if strings.HasPrefix(methodName, ENTITY_GET) {
			subUri := strings.ToLower(strings.TrimPrefix(methodName, ENTITY_GET))
			if subUri == "" {
				group.GET("/:id", fn.(echo.HandlerFunc))
				log.Debugf("Mapped GET /%s/:id", groupUri)
			}
		} else if strings.HasPrefix(methodName, ENTITY_LIST) {
			subUri := strings.ToLower(strings.TrimPrefix(methodName, ENTITY_LIST))
			group.GET(subUri, fn.(echo.HandlerFunc))
			log.Debugf("Mapped GET /%s", groupUri)
		} else if strings.HasPrefix(methodName, ENTITY_UPDATE) {
			subUri := strings.ToLower(strings.TrimPrefix(methodName, ENTITY_UPDATE))
			if subUri == "" {
				group.PUT("/:id", fn.(echo.HandlerFunc))
				log.Debugf("Mapped PUT /%s/:id", groupUri)
			}
		} else if strings.HasPrefix(methodName, ENTITY_DELETE) {
			subUri := strings.ToLower(strings.TrimPrefix(methodName, ENTITY_DELETE))
			if subUri == "" {
				group.DELETE("/:id", fn.(echo.HandlerFunc))
				log.Debugf("Mapped DELETE /%s/:id", groupUri)
			}
		} else {
			subUri := strings.ToLower(methodName)
			group.POST("/" + subUri, fn.(echo.HandlerFunc))
			log.Debugf("Mapped POST /%s/%s", groupUri, subUri)
		}
	}
}