package endpoint

import "github.com/labstack/echo"


type Endpoint interface {
	Create(o interface{}) echo.HandlerFunc
	Get(o interface{}) echo.HandlerFunc
	List(o interface{}) echo.HandlerFunc
	Update(o interface{}) echo.HandlerFunc
	Delete(o interface{}) echo.HandlerFunc
}