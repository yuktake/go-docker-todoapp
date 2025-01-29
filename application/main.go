package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type JwtCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwtv5.RegisteredClaims
}

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Content   string    `bun:"content,notnull"`
	Done      bool      `bun:"done"`
	Until     time.Time `bun:"until,nullzero"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type Data struct {
	Todos  []Todo
	Errors []error
}

// initEnvは.envファイルから環境変数を読み込みます
func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf(".envファイルの読み込みに失敗しました: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		return fmt.Errorf("JWT_SECRETが.envファイルに設定されていません")
	}

	return nil
}

// configureMiddlewareはEchoに必要なミドルウェアを設定します
func configureMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

func customFunc(todo *Todo) func([]string) []error {
	return func(values []string) []error {
		if len(values) == 0 || values[0] == "" {
			return nil
		}
		dt, err := time.Parse("2006-01-02T15:04 MST", values[0]+" JST")
		if err != nil {
			return []error{echo.NewBindingError("until", values[0:1], "failed to decode time", err)}
		}
		todo.Until = dt
		return nil
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func formatDateTime(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02 15:04")
}

func main() {
	// 環境変数の初期化
	if err := initEnv(); err != nil {
		log.Fatal(err)
	}

	sqldb, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(os.Getenv("DATABASE_URL"))
		log.Fatal(err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*Todo)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// 環境変数からJWTシークレットを取得
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwtv5.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: jwtSecret,
	})

	e.Renderer = &Template{
		templates: template.Must(template.New("").
			Funcs(template.FuncMap{
				"FormatDateTime": formatDateTime,
			}).ParseFS(templates, "templates/*")),
	}

	// パブリックルート: ログイン
	e.POST("/login", login)

	apiGroup := e.Group("/")
	apiGroup.Use(jwtMiddleware)

	apiGroup.GET("", func(c echo.Context) error {
		var todos []Todo
		ctx := context.Background()
		err := db.NewSelect().Model(&todos).Order("created_at").Scan(ctx)
		if err != nil {
			e.Logger.Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{
				Errors: []error{errors.New("Cannot get todos")},
			})
		}

		// JSONを返す
		return c.JSON(http.StatusOK, todos)
	})

	apiGroup.POST("", func(c echo.Context) error {
		var todo Todo
		// フォームパラメータをフィールドにバインド
		errs := echo.FormFieldBinder(c).
			Int64("id", &todo.ID).
			String("content", &todo.Content).
			Bool("done", &todo.Done).
			CustomFunc("until", customFunc(&todo)).
			BindErrors()
		if errs != nil {
			e.Logger.Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{Errors: errs})
		} else if todo.ID == 0 {
			// ID が 0 の時は登録
			ctx := context.Background()
			if todo.Content == "" {
				err = errors.New("Todo not found")
			} else {
				_, err = db.NewInsert().Model(&todo).Exec(ctx)
				if err != nil {
					e.Logger.Error(err)
					err = errors.New("Cannot update")
				}
			}
		} else {
			ctx := context.Background()
			if c.FormValue("delete") != "" {
				// 削除
				_, err = db.NewDelete().Model(&todo).Where("id = ?", todo.ID).Exec(ctx)
			} else {
				// 更新
				var orig Todo
				err = db.NewSelect().Model(&orig).Where("id = ?", todo.ID).Scan(ctx)
				if err == nil {
					orig.Done = todo.Done
					_, err = db.NewUpdate().Model(&orig).Where("id = ?", todo.ID).Exec(ctx)
				}
			}
			if err != nil {
				e.Logger.Error(err)
				err = errors.New("Cannot update")
			}
		}
		if err != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
		}
		return c.Redirect(http.StatusFound, "/")
	})

	staticFs, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FileSystem(http.FS(staticFs)))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fileServer)))
	e.Logger.Fatal(e.Start(":8000"))
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
