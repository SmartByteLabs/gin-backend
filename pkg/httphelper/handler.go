package httphelper

import (
	"encoding/json"
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func LoginHandler[T any, IDTYPE int64 | string](h rbac.UserHelper[T, IDTYPE], conditionFactory ConditionFactory[T]) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			NewResponse().Failed().AddError(err).SetMessage("Invalid request body").Send(http.StatusBadRequest, resp)
			return
		}

		condition, err := conditionFactory(req)
		if err != nil {
			NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.GetTableName()).Send(http.StatusBadRequest, resp)
			return
		}

		token, err := h.Login(req.Context(), body.Username, body.Password, condition)
		if err != nil {
			NewResponse().Failed().AddError(err).SetMessage("Error while login").Send(http.StatusInternalServerError, resp)
			return
		}

		NewResponse().Sucessfull().SetData(map[string]string{
			"token": token,
		}).Send(http.StatusOK, resp)
	}
}
