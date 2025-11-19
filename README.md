# Google Drive Go 封裝庫

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

一個簡潔易用的 Google Drive Golang 封裝庫。

## ✨ 主要特性

- 🔐 **Device Flow 授權** - 使用適用於"電視和受限輸入設備"的 OAuth2 授權模式
- 📤 **文件上傳** - 支持上傳文件到指定文件夾
- 🔄 **文件更新** - 按文件名稱更新覆蓋已存在的文件
- 🤖 **智能操作** - 自動判斷文件是否存在，不存在則創建，存在則更新
- 📁 **文件夾管理** - 支持創建和管理文件夾
- 🔑 **Token 自動刷新** - 自動處理 Token 過期和刷新
- 🌐 **瀏覽器引導** - 自動打開系統瀏覽器進行授權
- ⚡ **零依賴** - 純 Go 實現，無外部系統依賴

## 📦 安裝

```bash
go get github.com/Digman/gdrive
```

## 🚀 快速開始

### 1. 獲取 Google OAuth2 憑據

1. 訪問 [Google Cloud Console](https://console.cloud.google.com/)
2. 創建新項目或選擇現有項目
3. 啟用 Google Drive API
4. 創建 OAuth2 客戶端 ID（選擇"電視和受限輸入設備"類型）
5. 下載憑據文件，保存為 `credentials.json`

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
========================================
請完成以下步驟進行授權：
1. 系統將自動打開瀏覽器
2. 如果瀏覽器未自動打開，請手動訪問：https://...
3. 輸入以下代碼：XXXX-XXXX
4. 授權完成後，程序將自動繼續...
========================================
```

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

### Q: 如何獲取 OAuth2 憑據？
A: 訪問 [Google Cloud Console](https://console.cloud.google.com/)，創建"電視和受限輸入設備"類型的 OAuth2 客戶端 ID。

### Q: Token 過期了怎麼辦？
A: 程序會自動刷新 Token，無需手動處理。

### Q: 支持大文件上傳嗎？
A: 支持，但受 Google Drive API 限制。

### Q: 可以並發使用嗎？
A: 當前設計為單例使用，不建議並發操作。

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 許可證

MIT License

## 👤 作者

Digman

## 🙏 致謝

感謝 Google Drive API 團隊提供的優秀服務。
