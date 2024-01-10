package httphelper

import (
	"encoding/json"
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type ConditionFactory[T any] func(req *http.Request) (database.Condition[T], error)

type CrudHandler interface {
	GetModelName() string
	Get(id string, resp http.ResponseWriter, req *http.Request)
	GetAll(resp http.ResponseWriter, req *http.Request)
	Create(resp http.ResponseWriter, req *http.Request)
	Update(id string, resp http.ResponseWriter, req *http.Request)
	Delete(id string, resp http.ResponseWriter, req *http.Request)
}

type crudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string] struct {
	d        database.CrudHelper[T, MODEL, IDTYPE]
	idParser func(string) (IDTYPE, error)

	conditionFactory ConditionFactory[T]
}

func NewCrudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](d database.CrudHelper[T, MODEL, IDTYPE], idParser func(string) (IDTYPE, error), conditionFactory ConditionFactory[T]) CrudHandler {
	return crudHelper[T, MODEL, IDTYPE]{
		d:                d,
		idParser:         idParser,
		conditionFactory: conditionFactory,
	}
}

func (h crudHelper[T, MODEL, IDTYPE]) GetModelName() string {
	return h.d.GetTableName()
}

func (h crudHelper[T, MODEL, IDTYPE]) Get(idStr string, resp http.ResponseWriter, req *http.Request) {
	id, err := h.idParser(idStr)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid id for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while GET "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	if len(data) == 0 {
		NewResponse().Failed().SetMessage("No data found").Send(http.StatusNotFound, resp)
		return
	}

	NewResponse().Sucessfull().SetData(data[0]).Send(http.StatusOK, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) GetAll(resp http.ResponseWriter, req *http.Request) {
	condition, err := h.conditionFactory(req)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while GET "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	NewResponse().Sucessfull().SetData(data).Send(http.StatusOK, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Create(resp http.ResponseWriter, req *http.Request) {
	var body MODEL
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid data for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	data, err := h.d.Create(req.Context(), &body, condition)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while POST "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	NewResponse().Sucessfull().SetData(data).Send(http.StatusOK, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Update(idStr string, resp http.ResponseWriter, req *http.Request) {
	var body MODEL
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid data for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	id, err := h.idParser(idStr)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid id for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	err = h.d.Update(req.Context(), &body, req.URL.Query()["project"], condition.And(condition.New().Set("id", database.ConditionOperationEqual, id)))
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while PUT "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	data, err := h.d.Get(req.Context(), req.URL.Query()["project"], condition.And(condition.New().Set("id", database.ConditionOperationEqual, id)))
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while GET "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	if len(data) == 0 {
		NewResponse().Failed().SetMessage("No data found").Send(http.StatusNotFound, resp)
		return
	}

	NewResponse().Sucessfull().SetData(data[0]).Send(http.StatusOK, resp)
}

func (h crudHelper[T, MODEL, IDTYPE]) Delete(idStr string, resp http.ResponseWriter, req *http.Request) {
	id, err := h.idParser(idStr)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid id for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition, err := h.conditionFactory(req)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("In valid condition for "+h.d.GetTableName()).Send(http.StatusBadRequest, resp)
		return
	}

	condition.And(condition.New().Set("id", database.ConditionOperationEqual, id))

	err = h.d.Delete(req.Context(), condition)
	if err != nil {
		NewResponse().Failed().AddError(err).SetMessage("Error while DELETE "+h.d.GetTableName()).Send(http.StatusInternalServerError, resp)
		return
	}

	NewResponse().Sucessfull().Send(http.StatusOK, resp)
}
