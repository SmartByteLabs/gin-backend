package middleware

import (
	"net/http"
	"time"

	"github.com/princeparmar/gin-backend.git/pkg/config"
	"github.com/princeparmar/gin-backend.git/pkg/constant"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
	"github.com/princeparmar/gin-backend.git/pkg/utils"
)

type MiddlewareFuncWithNext func(http.ResponseWriter, *http.Request, func())

func CORSMiddleware(conf *config.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*") // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "*") // TODO: change this to specific domain
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func RequestIDMiddleware() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		utils.AddValueInRequestContext(req, constant.CtxKey_RequestID, utils.GenerateUUID())
	}
}

func LoggerMiddleware(log logger.Logger) MiddlewareFuncWithNext {
	return func(res http.ResponseWriter, req *http.Request, next func()) {
		start := time.Now()
		uid := req.Context().Value(constant.CtxKey_RequestID)
		log = log.WithField(constant.CtxKey_RequestID, uid)

		log.Infof("Request received from %v %v %v", req.RemoteAddr, req.Method, req.URL)
		utils.AddValueInRequestContext(req, constant.CtxKey_Logger, log)

		next()

		log.Infof("Request completed in %v", time.Since(start))
	}
}
