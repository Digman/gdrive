package gdrive

import (
	"os"
	"path/filepath"
	"time"
)

// BackupScheduler 備份調度器
type BackupScheduler struct {
	config          *Config
	client          *Client
	ticker          *time.Ticker
	stopChan        chan struct{}
	lastBackupTimes map[string]time.Time // 記錄每個文件的上次備份時間
	logger          Logger               // 日志實例
}

// NewBackupScheduler 創建備份調度器
func NewBackupScheduler(config *Config, client *Client) *BackupScheduler {
	logger := config.Logger
	if logger == nil {
		logger = newDefaultLogger()
	}

	return &BackupScheduler{
		config:          config,
		client:          client,
		lastBackupTimes: make(map[string]time.Time),
		logger:          logger,
	}
}

// Start 啟動調度器（異步運行）
func (s *BackupScheduler) Start() {
	s.ticker = time.NewTicker(s.config.BackupInterval)
	s.stopChan = make(chan struct{})

	// 異步執行定時任務
	go func() {
		// 啟動時立即執行一次
		s.runBackup()

		for {
			select {
			case <-s.ticker.C:
				s.runBackup()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()

	s.logger.Infof("[Backup] 定時備份已啟動，間隔: %v", s.config.BackupInterval)
}

// Stop 停止調度器
func (s *BackupScheduler) Stop() {
	if s.stopChan != nil {
		close(s.stopChan)
		s.logger.Infof("[Backup] 定時備份已停止")
	}
}

// runBackup 執行一次備份任務
func (s *BackupScheduler) runBackup() {
	s.logger.Infof("[Backup] 開始備份任務...")

	files, err := s.scanFiles()
	if err != nil {
		s.logger.Errorf("[Backup] 掃描文件失敗: %v", err)
		return
	}

	if len(files) == 0 {
		s.logger.Infof("[Backup] 沒有文件需要備份")
		return
	}

	successCount := 0
	failCount := 0

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			s.logger.Warningf("[Backup] 訪問文件失敗 %s: %v", file, err)
			failCount++
			continue // 單個文件失敗不影響其他
		}

		// 檢查是否需要備份
		if !s.shouldBackup(file, fileInfo) {
			continue
		}

		// 執行上傳
		_, isNew, err := s.client.UploadOrUpdateFile(file)
		if err != nil {
			s.logger.Errorf("[Backup] 備份失敗 %s: %v", file, err)
			failCount++
			continue // 單個文件失敗不影響其他
		}

		// 記錄備份時間
		s.lastBackupTimes[file] = fileInfo.ModTime()
		successCount++

		if isNew {
			s.logger.Infof("[Backup] 已創建: %s", file)
		} else {
			s.logger.Infof("[Backup] 已更新: %s", file)
		}
	}

	s.logger.Infof("[Backup] 備份完成 - 成功: %d, 失敗: %d", successCount, failCount)
}

// scanFiles 掃描需要備份的文件列表
func (s *BackupScheduler) scanFiles() ([]string, error) {
	var files []string

	for _, path := range s.config.BackupPaths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			s.logger.Warningf("[Backup] 訪問路徑失敗 %s: %v", path, err)
			continue // 單個路徑失敗不影響其他
		}

		if fileInfo.IsDir() {
			// 是目錄：遞歸掃描
			err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return nil // 忽略單個文件錯誤
				}

				// 跳過目錄本身
				if info.IsDir() {
					return nil
				}

				// 檢查排除規則
				if s.matchExclude(filePath) {
					return nil
				}

				files = append(files, filePath)
				return nil
			})
			if err != nil {
				s.logger.Warningf("[Backup] 掃描目錄失敗 %s: %v", path, err)
			}
		} else {
			// 是文件：直接添加
			if !s.matchExclude(path) {
				files = append(files, path)
			}
		}
	}

	return files, nil
}

// shouldBackup 判斷文件是否需要備份
func (s *BackupScheduler) shouldBackup(filePath string, fileInfo os.FileInfo) bool {
	// 全量模式：總是備份
	if s.config.BackupFullMode {
		return true
	}

	// 增量模式：檢查修改時間
	lastBackup, exists := s.lastBackupTimes[filePath]
	if !exists {
		return true // 首次備份
	}

	return fileInfo.ModTime().After(lastBackup) // 文件已修改
}

// matchExclude 檢查文件是否匹配排除規則
func (s *BackupScheduler) matchExclude(filePath string) bool {
	if len(s.config.BackupExcludes) == 0 {
		return false
	}

	fileName := filepath.Base(filePath)

	for _, pattern := range s.config.BackupExcludes {
		matched, err := filepath.Match(pattern, fileName)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}

	return false
}
