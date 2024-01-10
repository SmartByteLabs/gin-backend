package httphelper

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/constant"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type MiddlewareFuncWithNext func(http.ResponseWriter, *http.Request, func())

type MiddlewareFuncWithAbort MiddlewareFuncWithNext

func CORSMiddleware(origin, header string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", origin) // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", header) // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func LoggerMiddleware(log logger.Logger) MiddlewareFuncWithNext {
	return func(res http.ResponseWriter, req *http.Request, next func()) {
		start := time.Now()
		log = log.WithField(constant.CtxKey_RequestID, utils.GenerateUUID())

		log.Infof("Request received from %v %v %v", req.RemoteAddr, req.Method, req.URL)
		utils.AddValueInRequestContext(req, constant.CtxKey_Logger, log)

		next()

		log.Infof("Request completed in %v", time.Since(start))
	}
}

// RecoveryMiddleware recovers from panic and send error response
func RecoveryMiddleware(log logger.Logger) MiddlewareFuncWithNext {
	return func(res http.ResponseWriter, req *http.Request, next func()) {
		defer func() {
			if err := recover(); err != nil {
				NewResponse().Failed().SetMessage("Internal Server Error").Send(http.StatusInternalServerError, res)
				log.Error(fmt.Errorf("%v", err), "Panic recovered")
			}
		}()

		next()
	}
}

func DatabaseConnectionMiddleware(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		utils.AddValueInRequestContext(req, constant.CtxKey_DbConnection, db)
	}
}

func JWTAuthMiddleware[USER database.TableWithID[IDTYPE], IDTYPE int64 | string](loginRequired bool, secret string) MiddlewareFuncWithAbort {
	return func(res http.ResponseWriter, req *http.Request, abort func()) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" && !loginRequired {
			return
		}

		user, err := rbac.JWTAuthValidate[USER](authHeader, secret)
		if err != nil {
			NewResponse().Failed().SetMessage("authentication failed").AddError(err).Send(http.StatusUnauthorized, res)
			abort()
			return
		}

		utils.AddValueInRequestContext(req, constant.CtxKey_User, user)
	}
}
