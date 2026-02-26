package user

import (
	"errors"
	"net/http"
	"strconv"

	dtos "movie/system/internal/user/DTOs"
	"movie/system/pkg"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService Service
}

// NewUserHandler creates a UserHandler with the provided Service.
func NewUserHandler(service Service) *UserHandler {
	return &UserHandler{UserService: service}
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserCreateDTO
	if err := pkg.ReadJSON(w, r, &dto); err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	created, err := h.UserService.AddUser(r.Context(), dto)
	if errors.Is(err, ErrDuplicateEmail) {
		pkg.Error(ErrDuplicateEmailAPI, w, r)
		return
	} else if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(created, "User created successfully", w)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers(r.Context())
	if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}
	pkg.Ok(users, "", w)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		pkg.BadRequestWithCustomMessage(w, r, "invalid user ID", (*any)(nil))
		return
	}

	u, err := h.UserService.GetUserByID(r.Context(), uint(userId))
	if errors.Is(err, ErrUserNotFound) {
		pkg.NotFound(w, r, (*any)(nil))
		return
	} else if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(u, "", w)
}

func (h *UserHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		pkg.BadRequestWithCustomMessage(w, r, "invalid user ID", (*any)(nil))
		return
	}

	var dto dtos.ChangeRoleDTO
	if err := pkg.ReadJSON(w, r, &dto); err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	if err := h.UserService.ChangeRole(r.Context(), uint(userId), dto); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			pkg.NotFound(w, r, (*any)(nil))
			return
		}
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok((*any)(nil), "User role updated successfully", w)
}
