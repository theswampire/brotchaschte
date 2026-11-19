package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type UserDTO struct {
	Email string `json:"email"`
}

type GetUsersRes struct {
	Ok    bool      `json:"ok"`
	Error string    `json:"error"`
	Users []UserDTO `json:"users"`
}

func (s *Server) getUsersHandler(c *echo.Context) error {
	users_rows, err := s.auth.GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetUsersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die User konnten nicht geladen werden."})
	}
	users := make([]UserDTO, len(users_rows))
	for i, u := range users_rows {
		users[i] = UserDTO{Email: u.Email}
	}

	return c.JSON(http.StatusOK, &GetUsersRes{Ok: true, Users: users})
}

type CreateUserRes struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (s *Server) createUserHandler(c *echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	if err := s.auth.CreateUser(c.Request().Context(), email, password); err != nil {
		return c.JSON(http.StatusInternalServerError, &CreateUserRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Der Benutzer konnte nicht angelegt werden."})
	}
	return c.JSON(http.StatusOK, &CreateUserRes{Ok: true})
}

type UpdateUserPasswordRes struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (s *Server) updatePasswordHandler(c *echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	if err := s.auth.UpdateUserPassword(c.Request().Context(), email, password); err != nil {
		return c.JSON(http.StatusInternalServerError, &UpdateUserPasswordRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Das Passwort konnte nicht geändert werden."})
	}
	return c.JSON(http.StatusOK, &UpdateUserPasswordRes{Ok: true})

}

type DeleteUserRes struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (s *Server) deleteUserHandler(c *echo.Context) error {
	email := c.QueryParam("email")
	if err := s.auth.DeleteUser(c.Request().Context(), email); err != nil {
		return c.JSON(http.StatusInternalServerError, &DeleteUserRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Der Benutzer konnte nicht gelöscht werden."})
	}
	return c.JSON(http.StatusOK, &DeleteUserRes{Ok: true})

}
