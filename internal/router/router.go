package router

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/model"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type gitRouterHandler struct {
	log  logger.Logger
	conf *config.Config
	db   *sql.DB
}

func NewRouterHandler(log logger.Logger, conf *config.Config, db *sql.DB) *gitRouterHandler {
	return &gitRouterHandler{
		log:  log,
		conf: conf,
		db:   db,
	}
}

func (h *gitRouterHandler) RegisterRoute(r *gin.Engine) {
	v1 := r.Group("api/v1")

	authMiddleware := ginhelper.HttpHandlerToGinHandlerWithAbort(httphelper.JWTAuthMiddleware[rbac.User[int64]](true, h.conf.App.JWTSecret))

	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetTemplateHelper(h.db), utils.ParseInt, templateConditionFactory), authMiddleware)
	ginhelper.Register(v1, httphelper.NewCrudHelper(model.GetTempleHelper(h.db), utils.ParseInt, templeConditionFactory), authMiddleware)

	// register all rbac routes
	rbacRoutes(v1, h.db, h.conf)
}

func EmptyCondition(req *http.Request) (database.Condition[database.MysqlCondition], error) {
	return database.NewMysqlConditionHelper(), nil
}
