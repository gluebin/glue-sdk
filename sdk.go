package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Container interface {
	Get(id string) (component interface{}, found bool)
	Set(id string, component interface{}, overwrite bool) (actual interface{}, existed bool)
}

type Glue interface {
	Container() Container
	Http() Http
	Sql() Sql
	Config() *viper.Viper
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	return e.Code
}

var RequestMalformError = &Error{
	Code: "request_malform",
}

var UnknownError = &Error{
	Code:        "unknown_error",
	Description: "The module encounter an unknown error",
}

var NotFoundError = &Error{
	Code:        "not_found",
	Description: "There is no record found",
}

type Dto interface {
	Validate() error
}

type Http interface {
	JwtSign(c *gin.Context, data map[string]interface{}, expiresAt int64) (actual map[string]interface{}, token string, err error)
	JwtValidate(c *gin.Context, token string) (map[string]interface{}, bool)
	NewSessionCookie() *http.Cookie
	AddIdentityProvider(provider func(data map[string]interface{}))
	Router() *gin.RouterGroup
	Identity(*gin.Context) map[string]interface{}
	Authorizer() gin.HandlerFunc
	Authorized(func(*gin.Context)) func(*gin.Context)
	BindAndValidate(c *gin.Context, dto Dto) (valid bool)
	SendJson(c *gin.Context, data interface{}, err error)
}

type Sql interface {
	DB(databaseId ...string) (db *gorm.DB)
}

type Module interface {
	Plug(Glue)
	Run()
}
