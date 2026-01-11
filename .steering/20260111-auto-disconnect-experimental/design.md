# 設計: auto_disconnect 実験的機能への変更

## 実装アプローチ

### 変更対象ファイル

| ファイル | 変更内容 |
|---------|---------|
| `internal/generator/generator.go` | ループベースのexit処理、prompt_marker収集 |
| `internal/generator/generator_test.go` | テストケースの更新 |
| `CHANGELOG.md` | 変更履歴の追記 |
| `README.md` | 実験的機能の明記、制限事項の追加 |
| `README.en.md` | 英語版の更新 |
| `docs/functional-design.md` | 実験的機能であることを追記 |

### 生成されるTTLコード

**変更前:**
```ttl
; === Auto Disconnect ===
; Disconnect from step 2
sendln 'exit'
pause 1

:SUCCESS
closett
end
```

**変更後:**
```ttl
; === Auto Disconnect ===
; Loop until connection is closed (handles multi-hop SSH and shell sessions)
do
    flushrecv
    sendln 'exit'
    wait '$ ' '# '
loop while result > 0

:SUCCESS
closett
end
```

### コード設計

#### 1. prompt_marker収集関数

```go
func collectPromptMarkers(cfg *config.Config, route []*config.RouteStep) []string {
    seen := make(map[string]bool)
    var markers []string

    for _, step := range route {
        profile := cfg.Profiles[step.Profile]
        if profile != nil && profile.PromptMarker != "" {
            if !seen[profile.PromptMarker] {
                seen[profile.PromptMarker] = true
                markers = append(markers, profile.PromptMarker)
            }
        }
    }

    return markers
}
```

#### 2. auto_disconnect生成関数

```go
func generateAutoDisconnect(promptMarkers []string) string {
    var sb strings.Builder

    sb.WriteString("; === Auto Disconnect ===\n")
    sb.WriteString("; Loop until connection is closed\n")
    sb.WriteString("do\n")
    sb.WriteString("    flushrecv\n")
    sb.WriteString("    sendln 'exit'\n")

    // Build wait command with all prompt markers
    sb.WriteString("    wait")
    for _, marker := range promptMarkers {
        sb.WriteString(fmt.Sprintf(" '%s'", marker))
    }
    sb.WriteString("\n")

    sb.WriteString("loop while result > 0\n")
    sb.WriteString("\n")

    sb.WriteString(successTemplate)

    return sb.String()
}
```

### 動作フロー

```
┌─────────────────────────────────────────┐
│ do                                       │
│   ├─ flushrecv (バッファクリア)          │
│   ├─ sendln 'exit'                       │
│   └─ wait 'prompt1' 'prompt2' ...        │
│       ├─ マッチ → result > 0 → ループ継続│
│       └─ 接続終了 → result = 0 → 終了    │
└─────────────────────────────────────────┘
         ↓
    closett (Tera Term終了)
```

## 影響範囲

- `auto_disconnect: true` を指定した場合のみ影響
- `auto_disconnect: false`（デフォルト）の動作は変更なし
- YAMLスキーマの変更なし

## テスト方針

1. 既存の`TestGenerate_AutoDisconnect`を更新
   - ループ構造の存在確認
   - prompt_markerがwaitコマンドに含まれることを確認

2. 手動テスト（オプション）
   - Tera Termでの実際の動作確認
