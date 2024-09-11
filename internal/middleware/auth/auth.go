package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

type IAuthMiddleware interface {
	SessionMiddleware() echo.MiddlewareFunc
}

type authMiddleware struct {
	store *sessions.CookieStore
}

func NewAuthMiddleware(store *sessions.CookieStore) *authMiddleware {
	return &authMiddleware{store: store}
}

func (m *authMiddleware) SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// セッションをGorillaセッションストアから取得
			session, err := store.Get(c.Request(), "session")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Session not found")
			}

			// セッション内の"userId"を確認
			userId, ok := session.Values["userId"].(string)
			if !ok || userId == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
			}

			// 認証が成功した場合、次のハンドラーを実行
			return next(c)
		}
	}
}
