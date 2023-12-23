package httphelper

import (
	"net/http"
	"time"

	"github.com/princeparmar/gin-backend.git/pkg/constant"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
	"github.com/princeparmar/gin-backend.git/pkg/utils"
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

func JWTAuthMiddleware[USER any](loginRequired bool, secret string) MiddlewareFuncWithAbort {
	return func(res http.ResponseWriter, req *http.Request, abort func()) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" && !loginRequired {
			return
		}

		user, err := rbac.JWTAuthValidate[USER](authHeader, secret)
		if err != nil {
			NewResponse().Failed().SetMessage("").AddError(err).Send(http.StatusUnauthorized, res)
			abort()
			return
		}

		utils.AddValueInRequestContext(req, constant.CtxKey_User, user)
	}
}
