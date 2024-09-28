package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	s "github.com/core-go/search"
	"github.com/gorilla/mux"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

const InternalServerError = "Internal Server Error"

type UserHandler struct {
	service     service.UserService
	Validate    func(context.Context, interface{}) ([]core.ErrorMessage, error)
	Error       func(context.Context, string, ...map[string]interface{})
	Map         map[string]int
	ParamIndex  map[string]int
	FilterIndex int
}

func NewUserHandler(service service.UserService, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]core.ErrorMessage, error)) *UserHandler {
	userType := reflect.TypeOf(model.User{})
	_, jsonMap, _ := core.BuildMapField(userType)
	paramIndex, filterIndex := s.BuildAttributes(reflect.TypeOf(model.UserFilter{}))
	return &UserHandler{service: service, Validate: validate, Map: jsonMap, Error: logError, ParamIndex: paramIndex, FilterIndex: filterIndex}
}

func (h *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.All(r.Context())
	if err != nil {
		h.Error(r.Context(), fmt.Sprintf("Error: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, users)
}
func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := h.service.Load(r.Context(), id)
	if err != nil {
		h.Error(r.Context(), fmt.Sprintf("Error to get user '%s': %s", id, err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		JSON(w, http.StatusNotFound, nil)
	} else {
		JSON(w, http.StatusOK, user)
	}
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	errors, er2 := h.Validate(r.Context(), &user)
	if er2 != nil {
		h.Error(r.Context(), er2.Error(), MakeMap(user))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		JSON(w, http.StatusUnprocessableEntity, errors)
		return
	}
	res, er3 := h.service.Create(r.Context(), &user)
	if er3 != nil {
		h.Error(r.Context(), er3.Error(), MakeMap(user))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if res > 0 {
		JSON(w, http.StatusCreated, user)
	} else {
		JSON(w, http.StatusConflict, res)
	}
}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	errors, er2 := h.Validate(r.Context(), &user)
	if er2 != nil {
		h.Error(r.Context(), er2.Error(), MakeMap(user))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		JSON(w, http.StatusUnprocessableEntity, errors)
		return
	}
	res, er3 := h.service.Update(r.Context(), &user)
	if er3 != nil {
		h.Error(r.Context(), er3.Error(), MakeMap(user))
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	if res > 0 {
		JSON(w, http.StatusOK, user)
	} else if res == 0 {
		JSON(w, http.StatusNotFound, res)
	} else {
		JSON(w, http.StatusConflict, res)
	}
}
func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	var user model.User
	body, er1 := core.BuildMapAndStruct(r, &user)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	jsonUser, er2 := core.BodyToJsonMap(r, user, body, []string{"id"}, h.Map)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "method", "patch"))
	errors, er3 := h.Validate(r.Context(), &user)
	if er3 != nil {
		h.Error(r.Context(), er3.Error(), MakeMap(user))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		JSON(w, http.StatusUnprocessableEntity, errors)
		return
	}
	res, er4 := h.service.Patch(r.Context(), jsonUser)
	if er4 != nil {
		h.Error(r.Context(), er4.Error(), MakeMap(user))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if res > 0 {
		JSON(w, http.StatusOK, jsonUser)
	} else if res == 0 {
		JSON(w, http.StatusNotFound, res)
	} else {
		JSON(w, http.StatusConflict, res)
	}
}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	res, err := h.service.Delete(r.Context(), id)
	if err != nil {
		h.Error(r.Context(), fmt.Sprintf("Error to delete user '%s': %s", id, err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res > 0 {
		JSON(w, http.StatusOK, res)
	} else {
		JSON(w, http.StatusNotFound, res)
	}
}
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := model.UserFilter{Filter: &s.Filter{}}
	err := s.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := s.GetOffset(filter.Limit, filter.Page)
	users, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, &s.Result{List: &users, Total: total})
}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}
func MakeMap(res interface{}, opts ...string) map[string]interface{} {
	key := "request"
	if len(opts) > 0 && len(opts[0]) > 0 {
		key = opts[0]
	}
	m := make(map[string]interface{})
	b, err := json.Marshal(res)
	if err != nil {
		return m
	}
	m[key] = string(b)
	return m
}
