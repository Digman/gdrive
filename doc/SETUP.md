# Google Drive 憑據設置指南

## 獲取 OAuth2 憑據（credentials.json）

### 步驟 1: 訪問 Google Cloud Console

1. 訪問 [Google Cloud Console](https://console.cloud.google.com/)
2. 登錄您的 Google 帳戶

### 步驟 2: 創建或選擇項目

1. 點擊頂部的項目選擇器
2. 點擊「新建項目」或選擇現有項目
3. 輸入項目名稱（例如：「My Drive App」）
4. 點擊「創建」

### 步驟 3: 啟用 Google Drive API

1. 在左側菜單中，點擊「API 和服務」→「庫」
2. 搜索「Google Drive API」
3. 點擊「Google Drive API」
4. 點擊「啟用」按鈕

### 步驟 4: 創建 OAuth2 憑據

#### 重要：選擇正確的應用類型

1. 在左側菜單中，點擊「API 和服務」→「憑據」
2. 點擊「+ 創建憑據」→「OAuth 客戶端 ID」

3. **首次設置需要配置同意屏幕：**
   - 點擊「配置同意屏幕」
   - 選擇「外部」（或「內部」如果您是 Google Workspace 用戶）
   - 點擊「創建」
   - 填寫應用名稱（例如：「My Drive App」）
   - 填寫用戶支持電子郵件
   - 填寫開發者聯繫信息
   - 點擊「保存並繼續」
   - 在「範圍」頁面，點擊「保存並繼續」
   - 在「測試用戶」頁面（如果是外部應用），添加您的測試用戶郵箱
   - 點擊「保存並繼續」
   - 點擊「返回控制台」

4. **創建 OAuth2 客戶端 ID：**
   - 再次點擊「+ 創建憑據」→「OAuth 客戶端 ID」
   - **應用類型**：選擇「**電視和受限輸入設備**」
   - 輸入名稱（例如：「My Drive Client」）
   - 點擊「創建」

### 步驟 5: 下載憑據文件

1. 創建完成後，會顯示客戶端 ID 和客戶端密鑰
2. 點擊「下載 JSON」按鈕
3. 將下載的文件重命名為 `credentials.json`
4. 將文件放在您的項目根目錄

### 步驟 6: 驗證憑據文件格式

打開 `credentials.json` 文件，確保格式如下：

```json
{
  "installed": {
    "client_id": "YOUR_CLIENT_ID.apps.googleusercontent.com",
    "project_id": "your-project-id",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "YOUR_CLIENT_SECRET",
    "redirect_uris": ["http://localhost"]
  }
}
```

**重要提示：**
- 必須包含 `"installed"` 鍵
- 如果您的憑據文件包含 `"web"` 鍵而不是 `"installed"`，說明您選擇了錯誤的應用類型
- 請重新創建憑據，選擇「電視和受限輸入設備」或「桌面應用」類型

## 常見問題

### Q1: 我的憑據文件包含 "web" 而不是 "installed"？

**原因：** 您創建的是「Web 應用」類型的 OAuth2 客戶端。

**解決方法：**
1. 返回 Google Cloud Console
2. 刪除現有的 OAuth2 客戶端 ID
3. 重新創建，選擇「電視和受限輸入設備」或「桌面應用」類型

### Q2: 錯誤 "missing redirect URL in the client_credentials.json"？

**原因：** 使用了錯誤類型的憑據文件（Web 應用類型）。

**解決方法：** 參考 Q1 的解決方法。

### Q3: 授權時顯示「此應用未經驗證」？

**原因：** 您的應用尚未通過 Google 的驗證流程。

**解決方法：**
- 對於個人使用：點擊「高級」→「前往 [您的應用名稱]（不安全）」
- 對於生產環境：需要提交應用進行 Google 驗證

### Q4: 如何添加測試用戶？

**步驟：**
1. 在 Google Cloud Console 中，進入「API 和服務」→「OAuth 同意屏幕」
2. 滾動到「測試用戶」部分
3. 點擊「+ 添加用戶」
4. 輸入測試用戶的 Gmail 地址
5. 點擊「保存」

## 首次授權流程

當您運行程序時，會看到以下簡潔的授權提示：

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  🔐 Google Drive 設備授權
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  1. 瀏覽器將自動打開授權頁面
  2. 網址：https://www.google.com/device
  3. 輸入授權碼：XXXX-XXXX (授權碼會以青色高亮顯示)

  ⏳ 等待授權...
```

**完整操作步驟：**

1. **啟動程序** - 運行程序後會自動顯示授權界面
2. **瀏覽器自動打開** - 系統會嘗試自動打開瀏覽器
   - 如果未自動打開，手動訪問顯示的網址
3. **選擇 Google 帳戶** - 在瀏覽器中選擇要授權的帳戶
4. **處理安全警告**（如果出現）
   - 看到「此應用未經驗證」警告
   - 點擊「高級」→「前往 [您的應用名稱]（不安全）」
5. **輸入授權碼** - 輸入終端顯示的授權碼（會以青色高亮）
6. **確認授權** - 查看權限請求，點擊「允許」
7. **完成授權** - 看到成功消息後，程序自動繼續

**授權成功後：**
```
  ✅ 授權成功！
```

授權完成後，Token 會自動保存到 `token.json` 文件，下次運行時無需重新授權。

## 安全建議

1. **不要提交憑據文件到版本控制：**
   ```bash
   # 確保 .gitignore 包含以下內容
   credentials.json
   token.json
   ```

2. **保護 Token 文件：**
   - `token.json` 包含訪問您 Google Drive 的權限
   - 不要分享或公開此文件

3. **定期刷新憑據：**
   - 如果懷疑憑據洩露，立即在 Google Cloud Console 中撤銷
   - 創建新的 OAuth2 客戶端 ID

## 參考鏈接

- [Google Cloud Console](https://console.cloud.google.com/)
- [Google Drive API 文檔](https://developers.google.com/drive/api/v3/about-sdk)
- [OAuth2 設備授權流程](https://developers.google.com/identity/protocols/oauth2/limited-input-device)
