package gdrive

import (
	"fmt"
	"time"
)

// Config Google Drive 模塊配置
type Config struct {
	Enabled         bool   // 是否啟用
	FolderName      string // 文件夾名稱
	CredentialsFile string // 憑據文件路徑
	TokenFile       string // Token 文件路徑

	// 定時備份配置
	BackupEnabled  bool          // 是否啟用定時備份
	BackupInterval time.Duration // 備份間隔（如 30*time.Minute, time.Hour）
	BackupPaths    []string      // 要備份的文件/目錄路徑列表
	BackupExcludes []string      // 排除的文件模式（支持通配符，如 "*.tmp"）
	BackupFullMode bool          // true=全量備份，false=僅備份修改的文件
	Logger         Logger        // 日志實例（可選，nil 則使用默認實現）
}

// Validate 驗證配置有效性
func (c *Config) Validate() error {
	if !c.Enabled {
		return fmt.Errorf("Google Drive 模塊未啟用")
	}
	if c.CredentialsFile == "" {
		return fmt.Errorf("憑據文件路徑不能為空")
	}
	if c.TokenFile == "" {
		return fmt.Errorf("Token 文件路徑不能為空")
	}
	if c.FolderName == "" {
		return fmt.Errorf("文件夾名稱不能為空")
	}

	// 驗證備份配置
	if c.BackupEnabled {
		if c.BackupInterval <= 0 {
			return fmt.Errorf("BackupInterval 必須大於 0")
		}
		if len(c.BackupPaths) == 0 {
			return fmt.Errorf("BackupPaths 不能為空")
		}
	}

	return nil
}
