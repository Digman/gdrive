# Google Drive Go 封裝庫

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

一個簡潔易用的 Google Drive Golang 封裝庫。

## ✨ 主要特性

- 🔐 **Device Flow 授權** - 使用適用於"電視和受限輸入設備"的 OAuth2 授權模式
- 📤 **文件上傳** - 支持上傳文件到應用管理的文件夾
- 🔄 **文件更新** - 按文件名稱更新覆蓋已存在的文件
- 🤖 **智能操作** - 自動判斷文件是否存在，不存在則創建，存在則更新
- 📁 **文件夾管理** - 支持創建和管理應用專屬的文件夾
- 🔑 **Token 自動刷新** - 自動處理 Token 過期和刷新
- 🌐 **瀏覽器引導** - 自動打開系統瀏覽器進行授權
- 🔒 **安全隔離** - 僅訪問應用創建的文件，不影響用戶其他文件
- ⚡ **零依賴** - 純 Go 實現，無外部系統依賴

## 📦 安裝

```bash
go get github.com/Digman/gdrive
```

## 🚀 快速開始

### 1. 獲取 Google OAuth2 憑據

**⚠️ 重要：必須使用正確的憑據類型**

1. 訪問 [Google Cloud Console](https://console.cloud.google.com/)
2. 創建新項目或選擇現有項目
3. 啟用 Google Drive API
4. 創建 OAuth2 客戶端 ID
   - **應用類型**：必須選擇「**電視和受限輸入設備**」或「桌面應用」
   - ❌ 不要選擇「Web 應用」類型
5. 下載憑據文件，保存為 `credentials.json`

**憑據文件必須包含 `"installed"` 鍵，而不是 `"web"` 鍵。**

> 📖 詳細設置指南請參閱：[doc/SETUP.md](doc/SETUP.md)

### 2. 基本使用

```go
package main

import (
    "fmt"
    "log"

    "github.com/Digman/gdrive"
)

func main() {
    // 配置
    config := &gdrive.Config{
        Enabled:         true,
        FolderName:      "我的備份",
        CredentialsFile: "credentials.json",
        TokenFile:       "token.json",
    }

    // 創建客戶端
    client, err := gdrive.NewClient(config)
    if err != nil {
        log.Fatalf("創建客戶端失敗: %v", err)
    }

    // 智能上傳或更新文件（推薦）
    fileID, isNew, err := client.UploadOrUpdateFile("test.txt")
    if err != nil {
        log.Fatalf("操作失敗: %v", err)
    }

    if isNew {
        fmt.Printf("✅ 文件已創建，ID: %s\n", fileID)
    } else {
        fmt.Printf("✅ 文件已更新，ID: %s\n", fileID)
    }
}
```

### 3. 首次授權

首次運行時，程序會：

1. 自動打開系統瀏覽器
2. 顯示授權 URL 和設備代碼
3. 引導您完成授權流程
4. 自動保存 Token 到 `token.json`

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  🔐 Google Drive 設備授權
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  1. 瀏覽器將自動打開授權頁面
  2. 網址：https://www.google.com/device
  3. 輸入授權碼：XXXX-XXXX (授權碼會以青色高亮顯示)

  ⏳ 等待授權...

  ✅ 授權成功！
```

### 4. 重要說明：權限範圍

⚠️ **本庫僅訪問應用創建的文件和文件夾**

由於使用 Device Flow 授權模式，權限範圍受到限制：

- ✅ **可以做什麼**：在應用創建的文件夾（如 "我的備份"）中上傳、更新、刪除文件
- ❌ **不能做什麼**：訪問或修改用戶手動創建的文件和文件夾

這是一個**安全特性**，確保應用只能管理自己的文件，不會影響用戶 Drive 中的其他內容。

詳細說明請參閱：[API 文檔 - 授權範圍](doc/API.md#授權範圍)

## 📖 API 文檔

### 配置結構

```go
type Config struct {
    Enabled         bool   // 是否啟用
    FolderName      string // 文件夾名稱
    CredentialsFile string // 憑據文件路徑
    TokenFile       string // Token 文件路徑
}
```

### 主要方法

#### 創建客戶端

```go
client, err := gdrive.NewClient(config)
```

#### 上傳文件

```go
fileID, err := client.UploadFile("localfile.txt")
```

#### 更新文件

```go
fileID, err := client.UpdateFile("localfile.txt")
```

#### 智能上傳或更新（推薦）

```go
fileID, isNew, err := client.UploadOrUpdateFile("localfile.txt")
```

#### 創建文件夾

```go
folderID, err := client.CreateFolder("文件夾名稱", "父文件夾ID")
```

詳細 API 文檔請參閱 [doc/API.md](doc/API.md)

## 💡 使用示例

查看 [examples/main.go](examples/main.go) 獲取完整示例。

### 示例 1: 批量上傳文件

```go
files := []string{"file1.txt", "file2.txt", "file3.txt"}

for _, file := range files {
    fileID, isNew, err := client.UploadOrUpdateFile(file)
    if err != nil {
        log.Printf("處理 %s 失敗: %v", file, err)
        continue
    }

    status := "更新"
    if isNew {
        status = "創建"
    }
    fmt.Printf("✅ %s: %s (ID: %s)\n", status, file, fileID)
}
```

### 示例 2: 創建多級文件夾

```go
// 創建主文件夾
mainFolderID, _ := client.CreateFolder("備份", "")

// 創建子文件夾
subFolderID, _ := client.CreateFolder("2024", mainFolderID)
```

## 🛠️ 項目結構

```
gdrive/
├── README.md          # 項目說明
├── go.mod             # Go 模塊定義
├── config.go          # 配置結構
├── auth.go            # OAuth2 認證
├── client.go          # 客戶端封裝
├── file.go            # 文件操作
├── folder.go          # 文件夾操作
├── examples/          # 使用示例
│   └── main.go
└── doc/               # 文檔
    └── API.md         # API 文檔
```

## 🔒 安全注意事項

1. **不要提交憑據文件** - 將 `credentials.json` 和 `token.json` 添加到 `.gitignore`
2. **保護 Token 文件** - Token 具有訪問您 Google Drive 的權限
3. **使用環境變量** - 生產環境建議使用環境變量存儲敏感信息

```gitignore
credentials.json
token.json
```

## 📋 依賴項

- `google.golang.org/api/drive/v3` - Google Drive API v3
- `golang.org/x/oauth2` - OAuth2 認證庫
- `golang.org/x/oauth2/google` - Google OAuth2 實現

## ❓ 常見問題

### Q: 遇到 "unauthorized_client" 錯誤？
**A:** 刪除 `token.json` 文件，然後重新運行程序進行授權：
```bash
rm -f token.json
go run examples/main.go
```
詳細解決方案請參閱：[故障排除指南](doc/TROUBLESHOOTING.md)

### Q: 如何獲取 OAuth2 憑據？
A: 訪問 [Google Cloud Console](https://console.cloud.google.com/)，創建"電視和受限輸入設備"類型的 OAuth2 客戶端 ID。詳細步驟請參閱：[設置指南](doc/SETUP.md)

### Q: Token 過期了怎麼辦？
A: 程序會自動刷新 Token。如果自動刷新失敗，刪除 `token.json` 重新授權即可。

### Q: 支持大文件上傳嗎？
A: 支持，但受 Google Drive API 限制。

### Q: 可以並發使用嗎？
A: 當前設計為單例使用，不建議並發操作。

### Q: 遇到其他問題？
A: 請查閱 [故障排除指南](doc/TROUBLESHOOTING.md) 獲取詳細的解決方案。

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 許可證

MIT License

## 👤 作者

Digman

## 🙏 致謝

感謝 Google Drive API 團隊提供的優秀服務。
