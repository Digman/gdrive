package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Client Google Drive 客戶端
type Client struct {
	config    *Config
	service   *drive.Service
	folderID  string           // 緩存文件夾 ID
	scheduler *BackupScheduler // 備份調度器
}

// NewClient 創建新的 Google Drive 客戶端
func NewClient(config *Config) (*Client, error) {
	// 驗證配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 獲取 OAuth2 客戶端
	httpClient, err := getOAuth2Client(config)
	if err != nil {
		return nil, fmt.Errorf("認證失敗: %w", err)
	}

	// 創建 Drive Service
	service, err := drive.NewService(nil, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("創建 Drive Service 失敗: %w", err)
	}

	client := &Client{
		config:  config,
		service: service,
	}

	// 初始化時獲取或創建目標文件夾
	folderID, err := client.GetOrCreateFolder()
	if err != nil {
		return nil, fmt.Errorf("初始化文件夾失敗: %w", err)
	}
	client.folderID = folderID

	return client, nil
}

// GetFolderID 獲取當前使用的文件夾 ID
func (c *Client) GetFolderID() string {
	return c.folderID
}

// StartBackup 啟動定時備份（非阻塞，異步執行）
func (c *Client) StartBackup() error {
	if !c.config.BackupEnabled {
		return fmt.Errorf("備份未啟用，請在配置中設置 BackupEnabled = true")
	}

	if c.scheduler != nil {
		return fmt.Errorf("備份已在運行中")
	}

	// 創建並啟動調度器
	scheduler := NewBackupScheduler(c.config, c)

	c.scheduler = scheduler
	scheduler.Start() // 異步啟動
	return nil
}

// StopBackup 停止定時備份
func (c *Client) StopBackup() {
	if c.scheduler != nil {
		c.scheduler.Stop()
		c.scheduler = nil
	}
}
