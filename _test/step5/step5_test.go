package sta19_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
)

func TestStation19(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Errorf("環境変数の読み込みに失敗しました: %v", err)
		return
	}

	dbPath := "../../.sqlite3/todo.db"

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("データベースのクローズに失敗しました: %v", err)
			return
		}
	})

	r := router.NewRouter(todoDB)
	logChain := alice.New(middleware.GetOS, middleware.GetAccessLog)
	r.Handle("/healthz", logChain.Then(handler.NewHealthzHandler()))
	hTODO := handler.NewTODOHandler(service.NewTODOService(todoDB))
	r.Handle("/todos", logChain.Append(middleware.BasicAuth).Then(hTODO))
	hPanic := handler.NewPanicHandler()
	r.Handle("/do-panic", logChain.Append(middleware.Recovery).Then(hPanic))
	srv := httptest.NewServer(r)
	defer srv.Close()

	testcases := map[string]struct {
		Path               string
		UserID             string
		Password           string
		WantHTTPStatusCode int
	}{
		"Authentication is not required(1)": {
			Path:               "/healthz",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Authentication is not required(2)": {
			Path:               "/do-panic",
			WantHTTPStatusCode: http.StatusOK,
		},
		"UserID and Password are correct": {
			Path:               "/todos",
			UserID:             os.Getenv("BASIC_AUTH_USER_ID"),
			Password:           os.Getenv("BASIC_AUTH_PASSWORD"),
			WantHTTPStatusCode: http.StatusOK,
		},
		"Password is incorrect": {
			Path:               "/todos",
			UserID:             os.Getenv("BASIC_AUTH_USER_ID"),
			Password:           "DETARAME",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
		"UserID is incorrect": {
			Path:               "/todos",
			UserID:             "MACHIGAI",
			Password:           os.Getenv("BASIC_AUTH_PASSWORD"),
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
		"UserID and Password are incorrect": {
			Path:               "/todos",
			UserID:             "MACHIGAI",
			Password:           "DETARAME",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
		"UserID and Password are empty": {
			Path:               "/todos",
			UserID:             "",
			Password:           "",
			WantHTTPStatusCode: http.StatusUnauthorized,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, srv.URL+tc.Path, nil)
			if err != nil {
				t.Errorf("リクエストの作成に失敗しました: %v", err)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの送信に失敗しました: %v", err)
				return
			}
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("レスポンスのクローズに失敗しました: %v", err)
					return
				}
			})

			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}
		})
	}
}
