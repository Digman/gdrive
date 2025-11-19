# Google Drive Go 封裝庫 API 文檔

## 概述

這是一個 Google Drive 的 Golang 封裝庫，用於第三方調用。所有註釋和輸出都使用繁體中文。

## 主要特性

- ✅ 使用 Device Flow 授權模式（適用於"電視和受限輸入設備"類型項目）
- ✅ 支持文件上傳、更新和智能上傳或更新
- ✅ 支持創建文件夾
- ✅ 支持定時備份（異步、可配置間隔、支持全量/增量模式）
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

    // 定時備份配置
    BackupEnabled  bool          // 是否啟用定時備份
    BackupInterval time.Duration // 備份間隔（如 30*time.Minute, time.Hour）
    BackupPaths    []string      // 要備份的文件/目錄路徑列表
    BackupExcludes []string      // 排除的文件模式（支持通配符，如 "*.tmp"）
    BackupFullMode bool          // true=全量備份，false=僅備份修改的文件
    Logger         Logger        // 日志實例（可選，nil 則使用默認實現）
}
```

**備份配置字段說明：**

- `BackupEnabled`: 是否啟用定時備份功能
- `BackupInterval`: 備份執行間隔，使用 `time.Duration` 類型（如 `30*time.Minute`, `time.Hour`）
- `BackupPaths`: 需要備份的文件或目錄路徑列表（支持混合配置）
- `BackupExcludes`: 要排除的文件模式，支持通配符（如 `*.tmp`, `*.cache`）
- `BackupFullMode`:
  - `true`: 全量備份模式，每次備份所有文件
  - `false`: 增量備份模式，僅備份修改過的文件（基於文件修改時間）
- `Logger`: 日志實例，用於控制備份過程的日志輸出
  - `nil`: 使用默認實現（輸出到標準輸出）
  - 自定義實現：可集成到任何日志系統（logrus, zap 等）

#### 方法

##### Validate() error

驗證配置有效性。

**返回值：**
- `error`: 如果配置無效則返回錯誤信息

---

## 日志接口

### Logger

備份日志接口，允許用戶集成自定義日志系統。

#### 接口定義

```go
type Logger interface {
    // Infof 信息級別日志
    Infof(format string, v ...interface{})

    // Warningf 警告級別日志
    Warningf(format string, v ...interface{})

    // Errorf 錯誤級別日志
    Errorf(format string, v ...interface{})
}
```

#### 使用場景

**默認行為（不提供 Logger）：**
```go
config := &gdrive.Config{
    BackupEnabled: true,
    // ... 其他配置 ...
    // Logger 未設置，使用默認實現（輸出到標準輸出）
}
```

**集成標準庫 log：**
```go
type StdLogger struct {
    logger *log.Logger
}

func (l *StdLogger) Infof(format string, v ...interface{}) {
    l.logger.Printf("[INFO] " + format, v...)
}

func (l *StdLogger) Warningf(format string, v ...interface{}) {
    l.logger.Printf("[WARN] " + format, v...)
}

func (l *StdLogger) Errorf(format string, v ...interface{}) {
    l.logger.Printf("[ERROR] " + format, v...)
}

config := &gdrive.Config{
    BackupEnabled: true,
    Logger: &StdLogger{logger: log.New(os.Stdout, "", log.LstdFlags)},
}
```

**集成 logrus：**
```go
type LogrusAdapter struct {
    logger *logrus.Logger
}

func (l *LogrusAdapter) Infof(format string, v ...interface{}) {
    l.logger.Infof(format, v...)
}

func (l *LogrusAdapter) Warningf(format string, v ...interface{}) {
    l.logger.Warnf(format, v...)
}

func (l *LogrusAdapter) Errorf(format string, v ...interface{}) {
    l.logger.Errorf(format, v...)
}

config := &gdrive.Config{
    BackupEnabled: true,
    Logger: &LogrusAdapter{logger: logrus.New()},
}
```

**集成 zap：**
```go
type ZapAdapter struct {
    logger *zap.SugaredLogger
}

func (l *ZapAdapter) Infof(format string, v ...interface{}) {
    l.logger.Infof(format, v...)
}

func (l *ZapAdapter) Warningf(format string, v ...interface{}) {
    l.logger.Warnf(format, v...)
}

func (l *ZapAdapter) Errorf(format string, v ...interface{}) {
    l.logger.Errorf(format, v...)
}

config := &gdrive.Config{
    BackupEnabled: true,
    Logger: &ZapAdapter{logger: zap.S()},
}
```

#### 日志級別說明

- **Info**: 正常操作、成功信息
  - 定時備份已啟動
  - 備份任務開始/完成
  - 文件已創建/已更新

- **Warn**: 非致命錯誤、警告信息
  - 訪問文件失敗
  - 訪問路徑失敗
  - 掃描目錄失敗

- **Error**: 關鍵錯誤、失敗信息
  - 掃描文件失敗
  - 備份失敗

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

## 定時備份操作

### StartBackup

##### StartBackup() error

啟動定時備份（異步執行，非阻塞）。

**返回值：**
- `error`: 錯誤信息

**前置條件：**
- 配置中 `BackupEnabled` 必須為 `true`
- `BackupInterval` 必須大於 0
- `BackupPaths` 不能為空

**注意事項：**
- 異步執行，不會阻塞主程序
- 啟動時會立即執行一次備份
- 後續按照 `BackupInterval` 間隔自動執行
- 單個文件失敗不影響其他文件的備份
- 重複調用會返回錯誤

**示例：**
```go
config := &gdrive.Config{
    Enabled:         true,
    FolderName:      "我的備份",
    CredentialsFile: "credentials.json",
    TokenFile:       "token.json",

    // 定時備份配置
    BackupEnabled:  true,
    BackupInterval: 30 * time.Minute,
    BackupPaths:    []string{"./data", "./logs"},
    BackupExcludes: []string{"*.tmp", "*.cache"},
    BackupFullMode: false, // 增量備份
}

client, err := gdrive.NewClient(config)
if err != nil {
    log.Fatalf("創建客戶端失敗: %v", err)
}

// 啟動定時備份
if err := client.StartBackup(); err != nil {
    log.Fatalf("啟動備份失敗: %v", err)
}

// 程序繼續運行，備份在後台自動執行
```

---

### StopBackup

##### StopBackup()

停止定時備份。

**注意事項：**
- 安全停止備份調度器
- 如果備份未啟動，調用此方法無任何效果
- 建議在程序退出前調用以確保資源正確釋放

**示例：**
```go
// 程序退出前停止備份
defer client.StopBackup()

// 或者在需要時手動停止
client.StopBackup()
```

---

### 備份行為說明

**全量備份模式** (`BackupFullMode = true`)：
- 每次備份所有配置的文件和目錄
- 適用於文件數量較少或需要確保完整備份的場景

**增量備份模式** (`BackupFullMode = false`)：
- 僅備份修改過的文件（基於文件修改時間）
- 首次備份會備份所有文件
- 後續備份僅上傳自上次備份後修改的文件
- 適用於文件數量較多的場景，節省帶寬和時間

**文件掃描：**
- 支持指定單個文件：`BackupPaths: []string{"./config.json"}`
- 支持指定目錄：`BackupPaths: []string{"./data"}`（會遞歸掃描所有文件）
- 支持混合配置：`BackupPaths: []string{"./config.json", "./data", "./logs"}`

**排除規則：**
- 使用通配符模式（如 `*.tmp`, `*.cache`）
- 僅匹配文件名，不匹配路徑
- 多個排除規則會依次檢查

**錯誤處理：**
- 單個文件備份失敗不會中斷整個備份任務
- 失敗的文件會輸出錯誤信息但不會拋出異常
- 備份任務會繼續處理剩餘文件

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

- `https://www.googleapis.com/auth/drive.file` - 訪問應用創建和打開的文件

**權限範圍說明：**

由於使用 Device Flow 授權模式，只能使用受限的權限範圍。`drive.file` 範圍允許應用：

✅ **可以執行的操作：**
- 創建新文件夾
- 在應用創建的文件夾中上傳文件
- 更新應用創建的文件
- 刪除應用創建的文件和文件夾
- 管理應用創建的文件夾結構

❌ **無法執行的操作：**
- 訪問用戶手動創建的文件夾和文件
- 修改非應用創建的文件
- 列出整個 Drive 的所有文件

**這對您意味著什麼？**

當您使用本庫時：
1. 首次運行會在您的 Drive 中創建指定的文件夾（如 "我的備份"）
2. 所有上傳的文件都會保存在這個文件夾中
3. 應用只能看到和管理自己創建的文件和文件夾
4. 這是一個安全的設計，限制了應用的訪問範圍

---

## 注意事項

1. **憑據文件類型（重要）**：
   - ⚠️ 必須使用「**電視和受限輸入設備**」或「桌面應用」類型的 OAuth2 客戶端
   - ❌ 不要使用「Web 應用」類型
   - ✅ 憑據文件必須包含 `"installed"` 鍵，而不是 `"web"` 鍵
   - 📖 詳細設置指南：[SETUP.md](SETUP.md)

2. **Token 安全**：請妥善保管 `token.json` 文件，不要提交到版本控制系統

3. **文件大小**：默認不限制文件大小，但受 Google Drive API 限制

4. **並發使用**：當前設計為單例使用，不支持並發操作

5. **網絡要求**：需要能夠訪問 Google API 服務
