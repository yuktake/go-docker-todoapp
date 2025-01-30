package router

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/handler"
)

type Todo = todo.Todo

type Data struct {
	Todos  []Todo
	Errors []error
}

type JwtCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwtv5.RegisteredClaims
}

// ルーティング設定を行う関数
func RegisterRoutes(e *echo.Echo, todoHandler *handler.TodoHandler) {
	// 環境変数からJWTシークレットを取得
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwtv5.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: jwtSecret,
	})

	// パブリックルート: ログイン
	e.POST("/login", login)

	apiGroup := e.Group("/")
	apiGroup.Use(jwtMiddleware)

	apiGroup.GET("", todoHandler.GetTodos)

	apiGroup.POST("", todoHandler.CreateTodo)
}

// loginはユーザー認証を行い、JWTトークンを生成して返します
func login(c echo.Context) error {
	// リクエストからユーザー情報を取得
	email := c.FormValue("email")
	password := c.FormValue("password")

	// 簡易的な認証チェック（実際のアプリケーションではデータベースと連携）
	if email != "user@example.com" || password != "password" {
		log.Println("認証エラー:", email)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "メールアドレスまたはパスワードが無効です"})
	}

	// JWTクレームの設定
	claims := &JwtCustomClaims{
		Email: email,
		Name:  "User",
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}

	// クレームを持つトークンを生成
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	// トークンを署名
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("トークンの署名エラー:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "トークンを生成できませんでした"})
	}

	// トークンを返す
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// func getTodo(c echo.Context) error {
// 	var todos []Todo
// 	ctx := context.Background()
// 	err := database.NewSelect().Model(&todos).Order("created_at").Scan(ctx)
// 	if err != nil {
// 		e.Logger.Error(err)
// 		return c.Render(http.StatusBadRequest, "index", Data{
// 			Errors: []error{errors.New("Cannot get todos")},
// 		})
// 	}

// 	// JSONを返す
// 	return c.JSON(http.StatusOK, todos)
// }

// func createTodo(c echo.Context) error {
// 	var todo Todo
// 	// フォームパラメータをフィールドにバインド
// 	errs := echo.FormFieldBinder(c).
// 		Int64("id", &todo.ID).
// 		String("content", &todo.Content).
// 		Bool("done", &todo.Done).
// 		CustomFunc("until", customFunc(&todo)).
// 		BindErrors()
// 	if errs != nil {
// 		e.Logger.Error(err)
// 		return c.Render(http.StatusBadRequest, "index", Data{Errors: errs})
// 	} else if todo.ID == 0 {
// 		// ID が 0 の時は登録
// 		ctx := context.Background()
// 		if todo.Content == "" {
// 			err = errors.New("Todo not found")
// 		} else {
// 			_, err = database.NewInsert().Model(&todo).Exec(ctx)
// 			if err != nil {
// 				e.Logger.Error(err)
// 				err = errors.New("Cannot update")
// 			}
// 		}
// 	} else {
// 		ctx := context.Background()
// 		if c.FormValue("delete") != "" {
// 			// 削除
// 			_, err = database.NewDelete().Model(&todo).Where("id = ?", todo.ID).Exec(ctx)
// 		} else {
// 			// 更新
// 			var orig Todo
// 			err = database.NewSelect().Model(&orig).Where("id = ?", todo.ID).Scan(ctx)
// 			if err == nil {
// 				orig.Done = todo.Done
// 				_, err = database.NewUpdate().Model(&orig).Where("id = ?", todo.ID).Exec(ctx)
// 			}
// 		}
// 		if err != nil {
// 			e.Logger.Error(err)
// 			err = errors.New("Cannot update")
// 		}
// 	}
// 	if err != nil {
// 		return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
// 	}
// 	return c.Redirect(http.StatusFound, "/")
// }

// Fx Module
var Module = fx.Module("router",
	fx.Invoke(RegisterRoutes),
)
