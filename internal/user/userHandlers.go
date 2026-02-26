package user

import (
	"errors"
	dtos "movie/system/internal/user/DTOs"
	"movie/system/pkg"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService UserService
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var createDto dtos.UserCreateDTO
	errParse := pkg.ReadJSON(w, r, &createDto)
	if errParse != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	validateErrors := pkg.ValidateStruct(createDto)
	if validateErrors != nil {
		pkg.BadRequest(w, r, &validateErrors)
		return
	}

	user, errCreate := h.UserService.AddUser(r.Context(), createDto)
	if errors.Is(errCreate, ErrDuplicateEmail) {
		pkg.Error(ErrDuplicateEmailAPI, w, r)
		return
	} else if errCreate != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(user, "", w)
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
		pkg.BadRequestWithCustomMessage(w, r, err.Error(), (*any)(nil))
		return
	}

	user, err := h.UserService.GetUserByID(r.Context(), uint(userId))
	if errors.Is(err, ErrUserNotFound) {
		pkg.Error(ErrUserNotFoundAPI(err), w, r)
		return
	} else if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(user, "", w)
}

func (u *UserHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		pkg.BadRequestWithCustomMessage(w, r, err.Error(), (*any)(nil))
		return
	}
	var role dtos.ChangeRoleDTO

	errSeril := pkg.ReadJSON(w, r, &role)
	if errSeril != nil {
		pkg.BadRequestWithCustomMessage(w, r, errSeril.Error(), (*any)(nil))
		return
	}

	validateErrors := pkg.ValidateStruct(role)
	if validateErrors != nil {
		pkg.BadRequest(w, r, &validateErrors)
		return
	}

	updateErr := u.UserService.ChangeRole(r.Context(), uint(userId), role)
	if errors.Is(updateErr, ErrUserNotFound) {
		pkg.Error(ErrUserNotFoundAPI(updateErr), w, r)
		return
	}
	if updateErr != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok((*any)(nil), "User updated successfully", w)
}
