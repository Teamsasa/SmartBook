package usecaseUtils

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func CreateSession(c echo.Context, userId string) error {
	session, err := store.Get(c.Request(), "session")
	if err != nil {
		return err
	}

	session.Values["userId"] = userId

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}
