package gdrive

import "fmt"

// Config Google Drive 模塊配置
type Config struct {
	Enabled         bool   // 是否啟用
	FolderName      string // 文件夾名稱
	CredentialsFile string // 憑據文件路徑
	TokenFile       string // Token 文件路徑
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
	return nil
}
