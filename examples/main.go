package main

import (
	"fmt"
	"log"

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

	fmt.Println("\n🎉 所有操作完成！")
}
