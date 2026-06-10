package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
)

// loginHandler is the api handler for the login flow
func (s *Server) loginHandler(c *echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	ok, err := s.auth.AuthenticateUser(c.Request().Context(), email, password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "ein Fehler ist aufgetreten")
	}
	if !ok {
		return c.String(http.StatusUnauthorized, "invalid email or password")
	}

	sessionManager := s.auth.SessionManager()

	// Regenerate the session token to prevent session fixation
	if err := sessionManager.RenewToken(c.Request().Context()); err != nil {
		return c.String(http.StatusInternalServerError, "session error")
	}

	sessionManager.Put(c.Request().Context(), "email", email)
	return c.String(http.StatusOK, "logged in")
}

// logoutHandler is the api handler for the logout flow
func (s *Server) logoutHandler(c *echo.Context) error {
	sessionManager := s.auth.SessionManager()
	if err := sessionManager.Destroy(c.Request().Context()); err != nil {
		return c.String(http.StatusInternalServerError, "logout failed")
	}
	return c.NoContent(http.StatusOK)
}

// checkAuthenticatedHandler is the api handler to check if one is authenticated
func (s *Server) checkAuthenticatedHandler(c *echo.Context) error {
	// renew session
	sessionManager := s.auth.SessionManager()
	if err := sessionManager.RenewToken(c.Request().Context()); err != nil {
		return c.String(http.StatusInternalServerError, "session error")
	}
	// check if valid
	email := sessionManager.GetString(c.Request().Context(), "email")
	if email == "" {
		return c.String(http.StatusUnauthorized, "not authenticated")
	}
	return c.String(http.StatusOK, "authenticated")
}

// requireSession is a middleware to guard api endpoints. It check if a valid session exists.
func (s *Server) requireSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		email := s.auth.SessionManager().GetString(c.Request().Context(), "email")
		if email == "" {
			return c.String(http.StatusUnauthorized, "not authenticated")
		}
		return next(c)
	}
}

/// github.com/alexedwards/scs/v2 LoadAndSave middleware rewritten to use echo's context

// loadAndSaveHttpSession is the rewrite of the LoadAndSave middleware provided by scs. This version is compatible with echo.
// See source code of github.com/alexedwards/scs/v2 for more info
func (s *Server) loadAndSaveHttpSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sm := s.auth.SessionManager()

		c.Response().Header().Add("Vary", "Cookie")
		var token string
		cookie, err := c.Cookie(sm.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}

		ctx, err := sm.Load(c.Request().Context(), token)
		if err != nil {
			sm.ErrorFunc(c.Response(), c.Request(), err)
			return nil
		}

		sr := c.Request().WithContext(ctx)
		sw := &sessionResponseWriter{
			ResponseWriter: c.Response(),
			request:        sr,
			sessionManager: sm,
		}
		new_c := echo.NewContext(sr, sw)
		err = next(new_c)
		if !sw.written {
			commitAndWriteSessionCookie(sm, c.Response(), sr)
		}
		return err
	}

}

type sessionResponseWriter struct {
	http.ResponseWriter
	request        *http.Request
	sessionManager *scs.SessionManager
	written        bool
}

func (sw *sessionResponseWriter) Write(b []byte) (int, error) {
	if !sw.written {
		commitAndWriteSessionCookie(sw.sessionManager, sw.ResponseWriter, sw.request)
		sw.written = true
	}

	return sw.ResponseWriter.Write(b)
}

func (sw *sessionResponseWriter) WriteHeader(code int) {
	if !sw.written {
		commitAndWriteSessionCookie(sw.sessionManager, sw.ResponseWriter, sw.request)
		sw.written = true
	}

	sw.ResponseWriter.WriteHeader(code)
}

func (sw *sessionResponseWriter) Unwrap() http.ResponseWriter {
	return sw.ResponseWriter
}

func commitAndWriteSessionCookie(s *scs.SessionManager, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch s.Status(ctx) {
	case scs.Modified:
		token, expiry, err := s.Commit(ctx)
		if err != nil {
			s.ErrorFunc(w, r, err)
			return
		}

		s.WriteSessionCookie(ctx, w, token, expiry)
	case scs.Destroyed:
		s.WriteSessionCookie(ctx, w, "", time.Time{})
	}
}
