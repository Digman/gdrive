package gdrive

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

// UploadFile 上傳文件到配置的文件夾
// localPath: 本地文件路徑
// 返回: 文件 ID 和錯誤信息
func (c *Client) UploadFile(localPath string) (string, error) {
	// 打開本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("無法打開本地文件: %w", err)
	}
	defer file.Close()

	// 獲取文件名
	fileName := filepath.Base(localPath)

	// 創建 Drive 文件元數據
	driveFile := &drive.File{
		Name:    fileName,
		Parents: []string{c.folderID},
	}

	// 上傳文件
	createdFile, err := c.service.Files.Create(driveFile).
		Media(file).
		Fields("id, name").
		Do()
	if err != nil {
		return "", fmt.Errorf("上傳文件失敗: %w", err)
	}

	return createdFile.Id, nil
}

// UpdateFile 更新已存在的文件（按名稱查找並覆蓋）
// localPath: 本地文件路徑
// 返回: 文件 ID 和錯誤信息
func (c *Client) UpdateFile(localPath string) (string, error) {
	// 獲取文件名
	fileName := filepath.Base(localPath)

	// 查找已存在的文件
	fileID, err := c.findFileByName(fileName, c.folderID)
	if err != nil {
		return "", fmt.Errorf("查找文件失敗: %w", err)
	}

	// 打開本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("無法打開本地文件: %w", err)
	}
	defer file.Close()

	// 更新文件內容
	updatedFile, err := c.service.Files.Update(fileID, nil).
		Media(file).
		Fields("id, name").
		Do()
	if err != nil {
		return "", fmt.Errorf("更新文件失敗: %w", err)
	}

	return updatedFile.Id, nil
}

// UploadOrUpdateFile 智能上傳：不存在則創建，存在則更新
// localPath: 本地文件路徑
// 返回: 文件 ID、是否為新創建、錯誤信息
func (c *Client) UploadOrUpdateFile(localPath string) (string, bool, error) {
	// 獲取文件名
	fileName := filepath.Base(localPath)

	// 嘗試查找已存在的文件
	fileID, err := c.findFileByName(fileName, c.folderID)
	if err != nil {
		// 文件不存在，執行上傳
		fileID, err := c.UploadFile(localPath)
		if err != nil {
			return "", false, err
		}
		return fileID, true, nil
	}

	// 文件已存在，執行更新
	fileID, err = c.UpdateFile(localPath)
	if err != nil {
		return "", false, err
	}
	return fileID, false, nil
}

// findFileByName 根據文件名在指定文件夾中查找文件
// fileName: 文件名
// folderID: 文件夾 ID
// 返回: 文件 ID 和錯誤信息
func (c *Client) findFileByName(fileName, folderID string) (string, error) {
	// 構建查詢條件：文件名匹配、在指定文件夾中、未刪除
	query := fmt.Sprintf("name='%s' and '%s' in parents and trashed=false", fileName, folderID)

	// 執行查詢
	fileList, err := c.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		PageSize(1).
		Do()
	if err != nil {
		return "", fmt.Errorf("查詢文件失敗: %w", err)
	}

	// 檢查結果
	if len(fileList.Files) == 0 {
		return "", fmt.Errorf("文件不存在: %s", fileName)
	}

	return fileList.Files[0].Id, nil
}
