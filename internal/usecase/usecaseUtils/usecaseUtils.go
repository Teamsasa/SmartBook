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

	// セッションにユーザーIDを保存
	session.Values["userId"] = userId

	// Echoのレスポンスをhttp.ResponseWriterにキャストしてセッションを保存
	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return err
	}

	return nil
}
