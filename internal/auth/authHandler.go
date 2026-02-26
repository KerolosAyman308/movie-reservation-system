package auth

import (
	"errors"
	"net/http"
	"strconv"

	"movie/system/internal/user"
	dtos "movie/system/internal/user/DTOs"
	"movie/system/pkg"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	UserService   user.Service
	Authenticator Authenticator
	SecureCookie  bool
}

// NewAuthHandler creates an AuthHandler. Set secureCookie=true in production.
func NewAuthHandler(service user.Service, authenticator Authenticator, secureCookie bool) *AuthHandler {
	return &AuthHandler{
		UserService:   service,
		Authenticator: authenticator,
		SecureCookie:  secureCookie,
	}
}

type LoginResponse struct {
	User  dtos.UserResponse `json:"user"`
	Token string            `json:"token"`
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserSignUpDTO
	if err := pkg.ReadJSON(w, r, &dto); err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	defaultAdmin := false
	userToCreate := dtos.UserCreateDTO{
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Birthday:  dto.Birthday,
		IsAdmin:   &defaultAdmin,
		Password:  dto.Password,
	}

	createdUser, err := a.UserService.AddUser(r.Context(), userToCreate)
	if errors.Is(err, user.ErrDuplicateEmail) {
		pkg.Error(user.ErrDuplicateEmailAPI, w, r)
		return
	}
	if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	returnDto, err := a.generateToken(createdUser, w, r)
	if err != nil {
		return
	}
	pkg.Ok(returnDto, "Signed up successfully", w)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserLoginDTO
	if err := pkg.ReadJSON(w, r, &dto); err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	existUser, err := a.UserService.GetUserByEmail(r.Context(), dto.Email)
	if errors.Is(err, user.ErrUserNotFound) {
		pkg.Error(ErrNotAuthorized, w, r)
		return
	}
	if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	if ok := existUser.IsValidPassword(dto.Password); !ok {
		pkg.Error(ErrNotAuthorized, w, r)
		return
	}

	userDto := &dtos.UserResponse{
		ID:        existUser.ID,
		Email:     existUser.Email,
		FirstName: existUser.FirstName,
		LastName:  existUser.LastName,
		IsAdmin:   existUser.IsAdmin,
		Birthday:  existUser.Birthday,
	}

	returnDto, err := a.generateToken(userDto, w, r)
	if err != nil {
		return
	}
	pkg.Ok(returnDto, "Logged in successfully", w)
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		pkg.Unauthorized(w, r, (*any)(nil))
		return
	}

	token, err := a.Authenticator.ValidateRefreshToken(cookie.Value)
	if err != nil || !token.Valid {
		pkg.Unauthorized(w, r, (*any)(nil))
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		pkg.Error(ErrNotAuthorized, w, r)
		return
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	existUser, err := a.UserService.GetUserByID(r.Context(), uint(userID))
	if errors.Is(err, user.ErrUserNotFound) {
		pkg.Error(ErrNotAuthorized, w, r)
		return
	}
	if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	returnDto, err := a.generateToken(existUser, w, r)
	if err != nil {
		return
	}
	pkg.Ok(returnDto, "Token refreshed successfully", w)
}

func (a *AuthHandler) generateToken(existUser *dtos.UserResponse, w http.ResponseWriter, r *http.Request) (*LoginResponse, error) {
	accessToken, refreshToken, err := a.Authenticator.GenerateTokenPair(existUser.ID, existUser.IsAdmin)
	if err != nil {
		pkg.InternalError(w, r, (*any)(nil))
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   a.SecureCookie,
		Path:     "/api/v1/auth/refresh",
		MaxAge:   7 * 24 * 60 * 60,
		SameSite: http.SameSiteStrictMode,
	})

	return &LoginResponse{
		User: dtos.UserResponse{
			ID:        existUser.ID,
			Email:     existUser.Email,
			FirstName: existUser.FirstName,
			LastName:  existUser.LastName,
			Birthday:  existUser.Birthday,
			IsAdmin:   existUser.IsAdmin,
		},
		Token: accessToken,
	}, nil
}
