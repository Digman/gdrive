package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// getOAuth2Client 獲取已認證的 OAuth2 HTTP 客戶端
func getOAuth2Client(config *Config) (*http.Client, error) {
	ctx := context.Background()

	// 讀取憑據文件
	credentialsData, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("無法讀取憑據文件: %w", err)
	}

	// 解析 OAuth2 配置
	oauthConfig, err := google.ConfigFromJSON(credentialsData, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("無法解析憑據文件: %w", err)
	}

	// 嘗試從文件加載 Token
	token, err := loadToken(config.TokenFile)
	if err != nil {
		// Token 不存在或無效，需要重新認證
		token, err = getTokenFromDeviceFlow(ctx, oauthConfig)
		if err != nil {
			return nil, fmt.Errorf("設備認證失敗: %w", err)
		}

		// 保存 Token
		if err := saveToken(config.TokenFile, token); err != nil {
			return nil, fmt.Errorf("保存 Token 失敗: %w", err)
		}
	}

	// 創建 HTTP 客戶端（自動處理 Token 刷新）
	return oauthConfig.Client(ctx, token), nil
}

// getTokenFromDeviceFlow 通過 Device Flow 獲取新 Token
func getTokenFromDeviceFlow(ctx context.Context, oauthConfig *oauth2.Config) (*oauth2.Token, error) {
	// 獲取設備代碼
	deviceAuthResp, err := oauthConfig.DeviceAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("無法獲取設備代碼: %w", err)
	}

	// 顯示用戶授權信息
	fmt.Println("========================================")
	fmt.Println("請完成以下步驟進行授權：")
	fmt.Println("1. 系統將自動打開瀏覽器")
	fmt.Printf("2. 如果瀏覽器未自動打開，請手動訪問：%s\n", deviceAuthResp.VerificationURI)
	fmt.Printf("3. 輸入以下代碼：%s\n", deviceAuthResp.UserCode)
	fmt.Println("4. 授權完成後，程序將自動繼續...")
	fmt.Println("========================================")

	// 嘗試打開瀏覽器
	if err := openBrowser(deviceAuthResp.VerificationURI); err != nil {
		fmt.Printf("警告：無法自動打開瀏覽器: %v\n", err)
	}

	// 輪詢等待用戶授權
	token, err := oauthConfig.DeviceAccessToken(ctx, deviceAuthResp)
	if err != nil {
		return nil, fmt.Errorf("等待授權超時或失敗: %w", err)
	}

	fmt.Println("✅ 授權成功！")
	return token, nil
}

// saveToken 保存 Token 到文件
func saveToken(path string, token *oauth2.Token) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("無法創建 Token 文件: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(token); err != nil {
		return fmt.Errorf("無法寫入 Token: %w", err)
	}

	return nil
}

// loadToken 從文件加載 Token
func loadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("無法解析 Token 文件: %w", err)
	}

	// 檢查 Token 是否過期
	if token.Expiry.Before(time.Now()) && token.RefreshToken == "" {
		return nil, fmt.Errorf("Token 已過期且無法刷新")
	}

	return token, nil
}

// openBrowser 打開系統默認瀏覽器
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("不支持的操作系統: %s", runtime.GOOS)
	}

	return cmd.Start()
}
