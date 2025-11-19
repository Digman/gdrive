package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Digman/gdrive"
)

func main() {
	// 配置 Google Drive 客戶端
	config := &gdrive.Config{
		Enabled:         true,
		FolderName:      "我的備份",
		CredentialsFile: "credentials.json",
		TokenFile:       "token.json",
	}

	// 創建客戶端
	client, err := gdrive.NewClient(config)
	if err != nil {
		log.Fatalf("❌ 創建客戶端失敗: %v", err)
	}

	fmt.Println("✅ Google Drive 客戶端初始化成功")
	fmt.Printf("📁 使用文件夾 ID: %s\n\n", client.GetFolderID())

	// 示例 1: 上傳新文件
	fmt.Println("=== 示例 1: 上傳新文件 ===")
	fileID, err := client.UploadFile("test.txt")
	if err != nil {
		log.Printf("❌ 上傳失敗: %v\n", err)
	} else {
		fmt.Printf("✅ 文件已上傳，ID: %s\n\n", fileID)
	}

	// 示例 2: 更新已存在的文件
	fmt.Println("=== 示例 2: 更新已存在的文件 ===")
	fileID, err = client.UpdateFile("test.txt")
	if err != nil {
		log.Printf("❌ 更新失敗: %v\n", err)
	} else {
		fmt.Printf("✅ 文件已更新，ID: %s\n\n", fileID)
	}

	// 示例 3: 智能上傳或更新（推薦使用）
	fmt.Println("=== 示例 3: 智能上傳或更新 ===")
	fileID, isNew, err := client.UploadOrUpdateFile("test.txt")
	if err != nil {
		log.Fatalf("❌ 操作失敗: %v", err)
	}

	if isNew {
		fmt.Printf("✅ 文件已創建，ID: %s\n", fileID)
	} else {
		fmt.Printf("✅ 文件已更新，ID: %s\n", fileID)
	}

	// 示例 4: 創建新文件夾
	fmt.Println("\n=== 示例 4: 創建新文件夾 ===")
	newFolderID, err := client.CreateFolder("子文件夾", client.GetFolderID())
	if err != nil {
		log.Printf("❌ 創建文件夾失敗: %v\n", err)
	} else {
		fmt.Printf("✅ 文件夾已創建，ID: %s\n", newFolderID)
	}

	// 示例 5: 定時備份（可選）
	fmt.Println("\n=== 示例 5: 定時備份 ===")
	demonstrateBackup()

	// 示例 6: 自定義日志系統集成
	fmt.Println("\n=== 示例 6: 自定義日志系統集成 ===")
	demonstrateCustomLogger()

	fmt.Println("\n🎉 所有操作完成！")
}

// demonstrateBackup 演示定時備份功能
func demonstrateBackup() {
	// 配置定時備份
	config := &gdrive.Config{
		Enabled:         true,
		FolderName:      "我的備份",
		CredentialsFile: "credentials.json",
		TokenFile:       "token.json",

		// 定時備份配置
		BackupEnabled:  true,
		BackupInterval: 30 * time.Second,             // 演示用：每 30 秒備份一次（實際使用建議 30*time.Minute 或更長）
		BackupPaths:    []string{"./test_data"},      // 備份 test_data 目錄
		BackupExcludes: []string{"*.tmp", "*.cache"}, // 排除臨時文件
		BackupFullMode: false,                        // 增量備份：僅備份修改的文件
	}

	// 創建客戶端
	client, err := gdrive.NewClient(config)
	if err != nil {
		log.Printf("❌ 創建客戶端失敗: %v\n", err)
		return
	}

	// 啟動定時備份
	if err := client.StartBackup(); err != nil {
		log.Printf("❌ 啟動備份失敗: %v\n", err)
		return
	}

	fmt.Println("📝 定時備份示例運行中...")
	fmt.Println("   提示：實際使用時，主程序應保持運行以維持定時備份")
	fmt.Println("   示例將運行 2 分鐘後自動停止")

	// 運行 2 分鐘後停止（僅用於演示）
	time.Sleep(2 * time.Minute)

	// 停止備份
	client.StopBackup()
	fmt.Println("✅ 定時備份已停止")
}

// MyLogger 自定義日志實現（集成到用戶的日志系統）
type MyLogger struct{}

func (l *MyLogger) Infof(format string, v ...interface{}) {
	// 集成到自己的日志系統，例如：
	log.Printf("[INFO] "+format, v...)
}

func (l *MyLogger) Warningf(format string, v ...interface{}) {
	log.Printf("[WARN] "+format, v...)
}

func (l *MyLogger) Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// demonstrateCustomLogger 演示自定義日志系統集成
func demonstrateCustomLogger() {

	// 配置定時備份並注入自定義日志
	config := &gdrive.Config{
		Enabled:         true,
		FolderName:      "我的備份",
		CredentialsFile: "credentials.json",
		TokenFile:       "token.json",

		BackupEnabled:  true,
		BackupInterval: 30 * time.Second,
		BackupPaths:    []string{"./test_data"},
		Logger:         &MyLogger{}, // 使用自定義日志
	}

	client, err := gdrive.NewClient(config)
	if err != nil {
		log.Printf("❌ 創建客戶端失敗: %v\n", err)
		return
	}

	if err := client.StartBackup(); err != nil {
		log.Printf("❌ 啟動備份失敗: %v\n", err)
		return
	}

	fmt.Println("📝 使用自定義日志的備份運行中...")
	fmt.Println("   提示：備份日志將以 [INFO]/[WARN]/[ERROR] 前綴輸出")
	fmt.Println("   示例將運行 1 分鐘後自動停止")

	// 運行 1 分鐘後停止（僅用於演示）
	time.Sleep(1 * time.Minute)

	client.StopBackup()
	fmt.Println("✅ 自定義日志示例完成")
}
