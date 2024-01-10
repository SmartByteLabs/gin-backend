package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

func rbacRoutes(v1 *gin.RouterGroup, db *sql.DB, conf *config.Config) {
	//-------------------------------- Auth and Rbac related routes --------------------------------//
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewAccessHelper(db), utils.ParseInt, EmptyCondition))
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewRoleHelper(db), utils.ParseInt, EmptyCondition))
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewRoleAccessMappingHelper(db), utils.ParseInt, EmptyCondition))

	userHelper := rbac.NewUserHelper(mysql.NewUserHelper(db), conf.App.JWTSecret)
	v1.Group("auth").POST("/login", ginhelper.HttpHandlerToGinHandler(httphelper.LoginHandler(userHelper, EmptyCondition)))

	ginhelper.Register(v1, httphelper.NewCrudHelper(userHelper, utils.ParseInt, EmptyCondition))
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewUserRoleMappingHelper(db), utils.ParseInt, EmptyCondition))
	ginhelper.Register(v1, httphelper.NewCrudHelper(mysql.NewRequiredDataHelper(db), utils.ParseInt, EmptyCondition))
}
