package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

// CreateFolder 創建文件夾
// folderName: 文件夾名稱
// parentID: 父文件夾 ID（空字符串表示根目錄）
// 返回: 文件夾 ID 和錯誤信息
func (c *Client) CreateFolder(folderName, parentID string) (string, error) {
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	// 如果指定了父文件夾，則設置父級
	if parentID != "" {
		folder.Parents = []string{parentID}
	}

	// 創建文件夾
	createdFolder, err := c.service.Files.Create(folder).
		Fields("id, name").
		Do()
	if err != nil {
		return "", fmt.Errorf("創建文件夾失敗: %w", err)
	}

	return createdFolder.Id, nil
}

// GetOrCreateFolder 獲取或創建文件夾（不存在則創建）
// 返回: 文件夾 ID 和錯誤信息
func (c *Client) GetOrCreateFolder() (string, error) {
	folderName := c.config.FolderName

	// 先嘗試查找文件夾
	folderID, err := c.findFolderByName(folderName, "")
	if err == nil {
		// 文件夾已存在
		return folderID, nil
	}

	// 文件夾不存在，創建新文件夾
	folderID, err = c.CreateFolder(folderName, "")
	if err != nil {
		return "", fmt.Errorf("創建文件夾失敗: %w", err)
	}

	return folderID, nil
}

// findFolderByName 根據名稱查找文件夾
// folderName: 文件夾名稱
// parentID: 父文件夾 ID（空字符串表示在根目錄查找）
// 返回: 文件夾 ID 和錯誤信息
func (c *Client) findFolderByName(folderName, parentID string) (string, error) {
	// 構建查詢條件
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", folderName)
	if parentID != "" {
		query = fmt.Sprintf("%s and '%s' in parents", query, parentID)
	}

	// 執行查詢
	fileList, err := c.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		PageSize(1).
		Do()
	if err != nil {
		return "", fmt.Errorf("查詢文件夾失敗: %w", err)
	}

	// 檢查結果
	if len(fileList.Files) == 0 {
		return "", fmt.Errorf("文件夾不存在: %s", folderName)
	}

	return fileList.Files[0].Id, nil
}
