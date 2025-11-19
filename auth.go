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

// credentialsFile 憑據文件結構
type credentialsFile struct {
	Installed *struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"installed"`
}

// showCredentialsSetupGuide 顯示憑據設置指南並打開瀏覽器
func showCredentialsSetupGuide(credentialsPath string) {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  ⚠️  未找到或無法解析憑據文件")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("  預期路徑: %s\n", credentialsPath)
	fmt.Println()
	fmt.Println("  📝 請按照以下步驟獲取憑據文件：")
	fmt.Println()
	fmt.Println("  1. 訪問 Google Cloud Console（瀏覽器將自動打開）")
	fmt.Println("  2. 創建或選擇項目")
	fmt.Println("  3. 啟用 Google Drive API")
	fmt.Println("  4. 創建 OAuth2 憑據（類型：電視和受限輸入設備）")
	fmt.Println("  5. 下載憑據文件並保存為上述路徑")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 打開 Google Cloud Console
	consoleURL := "https://console.cloud.google.com/apis/credentials"
	if err := openBrowser(consoleURL); err != nil {
		fmt.Printf("  提示：無法自動打開瀏覽器，請手動訪問：\n  %s\n\n", consoleURL)
	} else {
		fmt.Println("  ✓ 已在瀏覽器中打開 Google Cloud Console")
		fmt.Println()
	}
}

// getOAuth2Client 獲取已認證的 OAuth2 HTTP 客戶端
func getOAuth2Client(config *Config) (*http.Client, error) {
	ctx := context.Background()

	// 讀取憑據文件
	credentialsData, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		showCredentialsSetupGuide(config.CredentialsFile)
		return nil, fmt.Errorf("無法讀取憑據文件: %w", err)
	}

	// 解析憑據文件
	var creds credentialsFile
	if err := json.Unmarshal(credentialsData, &creds); err != nil {
		showCredentialsSetupGuide(config.CredentialsFile)
		return nil, fmt.Errorf("無法解析憑據文件: %w", err)
	}

	if creds.Installed == nil {
		showCredentialsSetupGuide(config.CredentialsFile)
		return nil, fmt.Errorf("憑據文件格式錯誤：請使用「電視和受限輸入設備」或「已安裝應用」類型的 OAuth2 客戶端")
	}

	// 手動構建 OAuth2 配置（Device Flow）
	// 注意：Device Flow 不支持某些敏感權限範圍
	// 使用 drive.file 範圍，允許訪問應用創建和打開的文件
	oauthConfig := &oauth2.Config{
		ClientID:     creds.Installed.ClientID,
		ClientSecret: creds.Installed.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveFileScope},
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
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  🔐 Google Drive 設備授權")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  1. 瀏覽器將自動打開授權頁面")
	fmt.Printf("  2. 網址：%s\n", deviceAuthResp.VerificationURI)
	fmt.Printf("  3. 輸入授權碼：\033[1;36m%s\033[0m\n", deviceAuthResp.UserCode)
	fmt.Println()
	fmt.Println("  ⏳ 等待授權...")

	// 嘗試打開瀏覽器
	if err := openBrowser(deviceAuthResp.VerificationURI); err != nil {
		fmt.Printf("  ⚠️  無法自動打開瀏覽器，請手動訪問上方網址\n\n")
	}

	// 輪詢等待用戶授權
	token, err := oauthConfig.DeviceAccessToken(ctx, deviceAuthResp)
	if err != nil {
		return nil, fmt.Errorf("等待授權超時或失敗: %w", err)
	}

	fmt.Println("  ✅ 授權成功！")
	return token, nil
}

// saveToken 保存 Token 到文件
func saveToken(path string, token *oauth2.Token) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("無法創建 Token 文件: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(token); err != nil {
		_ = file.Close() // 忽略關閉錯誤，因為寫入已失敗
		return fmt.Errorf("無法寫入 Token: %w", err)
	}

	// 明確檢查 Close 錯誤
	if err := file.Close(); err != nil {
		return fmt.Errorf("無法關閉 Token 文件: %w", err)
	}

	return nil
}

// loadToken 從文件加載 Token
func loadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // 讀取操作可以忽略 Close 錯誤

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("無法解析 Token 文件: %w", err)
	}

	// 檢查 Token 是否過期
	if token.Expiry.Before(time.Now()) && token.RefreshToken == "" {
		return nil, fmt.Errorf("token 已過期且無法刷新")
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
