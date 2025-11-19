# 故障排除指南

## 常見錯誤及解決方案

### 錯誤 1: "invalid_scope" "Invalid device flow scope"

**完整錯誤信息：**
```
oauth2: "invalid_scope" "Invalid device flow scope: https://www.googleapis.com/auth/drive"
```

**原因：**
Device Flow 授權模式不支持完整的 Google Drive 訪問權限（`drive` 範圍）。

**解決方案：**

這個錯誤表明代碼已經過時。請確保您使用的是最新版本的代碼。

**檢查當前權限範圍：**

查看 `auth.go` 文件中的配置：

```bash
grep -A 5 "Scopes:" auth.go
```

應該顯示：
```go
Scopes: []string{drive.DriveFileScope},
```

如果顯示的是 `drive.DriveScope`，請更新為 `drive.DriveFileScope`。

**說明：**
- `drive.DriveFileScope` = `https://www.googleapis.com/auth/drive.file` ✅ Device Flow 支持
- `drive.DriveScope` = `https://www.googleapis.com/auth/drive` ❌ Device Flow 不支持

---

### 錯誤 2: "unauthorized_client" "Unauthorized"

**完整錯誤信息：**
```
oauth2: "unauthorized_client" "Unauthorized"
```

**原因：**
1. 使用了舊的授權 Token，授權範圍不正確
2. OAuth 同意屏幕配置不正確

**解決方案：**

#### 步驟 1: 刪除舊的 Token 文件

```bash
# 在項目根目錄或 examples 目錄下執行
rm -f token.json
```

#### 步驟 2: 檢查 OAuth 同意屏幕配置

1. 訪問 [Google Cloud Console](https://console.cloud.google.com/)
2. 選擇您的項目
3. 進入「API 和服務」→「OAuth 同意屏幕」
4. 確認「範圍」部分包含 Google Drive API 權限
5. 如果沒有，點擊「添加或移除範圍」
6. 搜索並添加：
   - `https://www.googleapis.com/auth/drive` - **推薦**
   - 或 `https://www.googleapis.com/auth/drive.file` - 僅訪問應用創建的文件

#### 步驟 3: 重新運行程序

```bash
go run examples/main.go
```

程序會自動引導您完成新的授權流程。

---

### 錯誤 2: "missing redirect URL in the client_credentials.json"

**完整錯誤信息：**
```
oauth2/google: missing redirect URL in the client_credentials.json
```

**原因：**
使用了錯誤類型的 OAuth2 憑據（Web 應用類型）。

**解決方案：**

1. 刪除當前的 `credentials.json`
2. 在 Google Cloud Console 中刪除現有的 OAuth2 客戶端 ID
3. 創建新的 OAuth2 客戶端 ID
   - **應用類型**：選擇「**電視和受限輸入設備**」或「桌面應用」
   - ❌ 不要選擇「Web 應用」
4. 下載新的憑據文件並保存為 `credentials.json`
5. 驗證文件格式包含 `"installed"` 鍵

詳細步驟請參閱：[SETUP.md](SETUP.md)

---

### 錯誤 3: "憑據文件格式錯誤：請使用「電視和受限輸入設備」或「已安裝應用」類型的 OAuth2 客戶端"

**原因：**
憑據文件不包含 `"installed"` 鍵。

**解決方案：**

打開 `credentials.json`，檢查格式：

**✅ 正確的格式：**
```json
{
  "installed": {
    "client_id": "...",
    "client_secret": "...",
    ...
  }
}
```

**❌ 錯誤的格式：**
```json
{
  "web": {
    "client_id": "...",
    ...
  }
}
```

如果是錯誤格式，請參考「錯誤 2」的解決方案重新創建憑據。

---

### 錯誤 4: "此應用未經驗證"

**授權時出現的警告：**
```
Google 尚未驗證此應用
```

**原因：**
您的應用尚未通過 Google 的驗證流程。

**解決方案（個人使用）：**

1. 點擊「高級」
2. 點擊「前往 [您的應用名稱]（不安全）」
3. 繼續授權流程

**解決方案（生產環境）：**

需要提交應用進行 Google 驗證：
1. 在 Google Cloud Console 中完善應用信息
2. 提交驗證申請
3. 等待 Google 審核（可能需要數周時間）

---

### 錯誤 5: Token 過期

**錯誤信息：**
```
Token 已過期且無法刷新
```

**解決方案：**

```bash
# 刪除過期的 Token
rm -f token.json

# 重新運行程序進行授權
go run examples/main.go
```

---

### 錯誤 6: "文件不存在" 當調用 UpdateFile 時

**錯誤信息：**
```
查找文件失敗: 文件不存在: filename.txt
```

**原因：**
嘗試更新一個不存在的文件。

**解決方案：**

使用 `UploadOrUpdateFile` 方法代替 `UpdateFile`：

```go
// ❌ 不推薦：文件可能不存在
fileID, err := client.UpdateFile("test.txt")

// ✅ 推薦：自動判斷創建或更新
fileID, isNew, err := client.UploadOrUpdateFile("test.txt")
```

---

## 調試技巧

### 1. 查看詳細錯誤信息

在代碼中添加詳細的錯誤日誌：

```go
client, err := gdrive.NewClient(config)
if err != nil {
    log.Printf("完整錯誤: %+v", err)
    return
}
```

### 2. 驗證憑據文件

```bash
# 查看憑據文件內容
cat credentials.json | jq .

# 檢查是否包含 "installed" 鍵
cat credentials.json | jq '.installed'
```

### 3. 檢查 Token 內容

```bash
# 查看 Token 文件內容
cat token.json | jq .

# 檢查過期時間
cat token.json | jq '.expiry'
```

### 4. 測試網絡連接

```bash
# 測試是否能訪問 Google API
curl -I https://www.googleapis.com/drive/v3/about
```

---

## 重置所有設置

如果遇到無法解決的問題，可以完全重置：

```bash
# 1. 刪除所有認證文件
rm -f credentials.json token.json

# 2. 在 Google Cloud Console 中：
#    - 撤銷應用的授權
#    - 刪除 OAuth2 客戶端 ID
#    - 重新創建（選擇正確的類型）

# 3. 重新開始設置流程
# 參考 SETUP.md
```

---

## 獲取幫助

如果上述方案都無法解決問題：

1. **檢查 Google Drive API 狀態**：
   - [Google Workspace Status Dashboard](https://www.google.com/appsstatus/dashboard/)

2. **查看 Google Drive API 文檔**：
   - [Google Drive API v3 文檔](https://developers.google.com/drive/api/v3/reference)

3. **提交 Issue**：
   - 包含完整的錯誤信息
   - 憑據文件類型（不要包含實際的 client_id 和 secret）
   - 您的操作步驟
