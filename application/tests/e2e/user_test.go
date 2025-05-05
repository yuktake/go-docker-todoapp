package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuktake/todo-webapp/dto"
)

func TestCreateUserSuccess(t *testing.T) {
	// 1. 他のテストで使用されたデータを削除
	if err := CleanupAllTables(testDB); err != nil {
		t.Fatalf("failed to cleanup tables: %v", err)
	}

	params := map[string]string{
		"name":     "testuser",
		"password": "testpassword",
		"email":    "test@test100.com",
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}

	// 2. HTTPリクエストを作成
	req, err := http.NewRequest(
		http.MethodPost,
		TestServer.URL+"/signup",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. HTTPクライアントを作成してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 4. レスポンスの検証
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", resp.StatusCode)
	}

	// 5. レスポンスボディ取得
	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var signupResp dto.SignupResponse
	if err := json.Unmarshal(body, &signupResp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, params["name"], signupResp.User.Name)
}
