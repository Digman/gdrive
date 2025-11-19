# Google Drive Go 封裝庫 API 文檔

## 概述

這是一個 Google Drive 的 Golang 封裝庫，用於第三方調用。所有註釋和輸出都使用繁體中文。

## 主要特性

- ✅ 使用 Device Flow 授權模式（適用於"電視和受限輸入設備"類型項目）
- ✅ 支持文件上傳、更新和智能上傳或更新
- ✅ 支持創建文件夾
- ✅ Token 自動刷新機制
- ✅ 系統瀏覽器自動引導授權
- ✅ 純 Go 實現，無外部依賴

---

## 配置結構

### Config

```go
type Config struct {
    Enabled         bool   // 是否啟用
    FolderName      string // 文件夾名稱
    CredentialsFile string // 憑據文件路徑
    TokenFile       string // Token 文件路徑
}
```

#### 方法

##### Validate() error

驗證配置有效性。

**返回值：**
- `error`: 如果配置無效則返回錯誤信息

---

## 客戶端

### Client

Google Drive 客戶端封裝。

#### 創建客戶端

##### NewClient(config *Config) (*Client, error)

創建新的 Google Drive 客戶端。

**參數：**
- `config`: 配置對象

**返回值：**
- `*Client`: 客戶端實例
- `error`: 錯誤信息

**示例：**
```go
config := &gdrive.Config{
    Enabled:         true,
    FolderName:      "我的備份",
    CredentialsFile: "credentials.json",
    TokenFile:       "token.json",
}

client, err := gdrive.NewClient(config)
if err != nil {
    log.Fatalf("創建客戶端失敗: %v", err)
}
```

#### 方法

##### GetFolderID() string

獲取當前使用的文件夾 ID。

**返回值：**
- `string`: 文件夾 ID

---

## 文件操作

### UploadFile

##### UploadFile(localPath string) (string, error)

上傳文件到配置的文件夾。

**參數：**
- `localPath`: 本地文件路徑

**返回值：**
- `string`: 文件 ID
- `error`: 錯誤信息

**示例：**
```go
fileID, err := client.UploadFile("test.txt")
if err != nil {
    log.Fatalf("上傳失敗: %v", err)
}
fmt.Printf("文件已上傳，ID: %s\n", fileID)
```

---

### UpdateFile

##### UpdateFile(localPath string) (string, error)

更新已存在的文件（按名稱查找並覆蓋）。

**參數：**
- `localPath`: 本地文件路徑

**返回值：**
- `string`: 文件 ID
- `error`: 錯誤信息

**注意事項：**
- 文件必須已存在於目標文件夾中
- 按文件名稱匹配
- 直接覆蓋文件內容

**示例：**
```go
fileID, err := client.UpdateFile("test.txt")
if err != nil {
    log.Fatalf("更新失敗: %v", err)
}
fmt.Printf("文件已更新，ID: %s\n", fileID)
```

---

### UploadOrUpdateFile（推薦）

##### UploadOrUpdateFile(localPath string) (string, bool, error)

智能上傳：文件不存在則創建，存在則更新。

**參數：**
- `localPath`: 本地文件路徑

**返回值：**
- `string`: 文件 ID
- `bool`: 是否為新創建（`true` 表示新創建，`false` 表示更新）
- `error`: 錯誤信息

**示例：**
```go
fileID, isNew, err := client.UploadOrUpdateFile("test.txt")
if err != nil {
    log.Fatalf("操作失敗: %v", err)
}

if isNew {
    fmt.Printf("文件已創建，ID: %s\n", fileID)
} else {
    fmt.Printf("文件已更新，ID: %s\n", fileID)
}
```

---

## 文件夾操作

### CreateFolder

##### CreateFolder(folderName, parentID string) (string, error)

創建文件夾。

**參數：**
- `folderName`: 文件夾名稱
- `parentID`: 父文件夾 ID（空字符串表示根目錄）

**返回值：**
- `string`: 文件夾 ID
- `error`: 錯誤信息

**示例：**
```go
// 在根目錄創建文件夾
folderID, err := client.CreateFolder("新文件夾", "")

// 在指定文件夾下創建子文件夾
subFolderID, err := client.CreateFolder("子文件夾", folderID)
```

---

### GetOrCreateFolder

##### GetOrCreateFolder() (string, error)

獲取或創建文件夾（不存在則創建）。

**返回值：**
- `string`: 文件夾 ID
- `error`: 錯誤信息

**注意事項：**
- 使用配置中的 `FolderName`
- 如果文件夾已存在則返回現有 ID
- 如果不存在則自動創建

---

## 授權流程

### Device Flow 授權

首次使用時，程序會自動啟動 Device Flow 授權：

1. **自動打開瀏覽器**：程序會嘗試打開系統默認瀏覽器
2. **手動訪問（如果瀏覽器未自動打開）**：訪問顯示的 URL
3. **輸入授權碼**：在瀏覽器中輸入顯示的授權碼
4. **完成授權**：授權成功後，程序自動繼續執行
5. **Token 持久化**：Token 會自動保存到配置的 `TokenFile`

### Token 自動刷新

- Token 會自動保存到指定文件
- 程序會自動檢測 Token 是否過期
- 過期的 Token 會自動刷新（如果有 RefreshToken）
- 無需手動處理 Token 刷新邏輯

---

## 錯誤處理

所有公開方法都返回 `error` 類型的錯誤信息：

```go
fileID, err := client.UploadFile("test.txt")
if err != nil {
    // 處理錯誤
    log.Printf("上傳失敗: %v", err)
    return
}
```

### 常見錯誤

| 錯誤信息 | 原因 | 解決方法 |
|---------|------|---------|
| `Google Drive 模塊未啟用` | `Config.Enabled` 為 `false` | 設置為 `true` |
| `憑據文件路徑不能為空` | `CredentialsFile` 未設置 | 提供有效的憑據文件路徑 |
| `Token 文件路徑不能為空` | `TokenFile` 未設置 | 提供有效的 Token 文件路徑 |
| `無法讀取憑據文件` | 憑據文件不存在或無權限 | 檢查文件路徑和權限 |
| `文件不存在` | 調用 `UpdateFile` 但文件不存在 | 使用 `UploadOrUpdateFile` 代替 |
| `設備認證失敗` | 授權過程中斷或超時 | 重新運行程序並完成授權 |

---

## 完整示例

參見 `examples/main.go` 文件。

```go
package main

import (
    "fmt"
    "log"
    "github.com/longjinhua/gdrive"
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

    // 智能上傳或更新
    fileID, isNew, err := client.UploadOrUpdateFile("test.txt")
    if err != nil {
        log.Fatalf("操作失敗: %v", err)
    }

    if isNew {
        fmt.Printf("文件已創建，ID: %s\n", fileID)
    } else {
        fmt.Printf("文件已更新，ID: %s\n", fileID)
    }
}
```

---

## 依賴項

- `google.golang.org/api/drive/v3` - Google Drive API v3
- `golang.org/x/oauth2` - OAuth2 認證庫
- `golang.org/x/oauth2/google` - Google OAuth2 實現

---

## 授權範圍

本庫使用以下 OAuth2 授權範圍：

- `https://www.googleapis.com/auth/drive.file` - 管理用戶 Google Drive 中的文件

---

## 注意事項

1. **憑據文件**：需要在 Google Cloud Console 創建"電視和受限輸入設備"類型的 OAuth2 客戶端，並下載憑據文件
2. **Token 安全**：請妥善保管 `token.json` 文件，不要提交到版本控制系統
3. **文件大小**：默認不限制文件大小，但受 Google Drive API 限制
4. **並發使用**：當前設計為單例使用，不支持並發操作
5. **網絡要求**：需要能夠訪問 Google API 服務
