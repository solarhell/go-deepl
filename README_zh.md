# go-deepl

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/solarhell/go-deepl?style=for-the-badge)](https://goreportcard.com/report/github.com/solarhell/go-deepl)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/solarhell/go-deepl)
[![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/solarhell/go-deepl?utm_source=oss&utm_medium=github&utm_campaign=solarhell%2Fgo-deepl&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)](https://coderabbit.ai)

[English](README.md) | 中文

轻量、零依赖的 [DeepL API](https://developers.deepl.com/docs) Go 客户端，基于 Go 1.26+。

## 安装

```bash
go get github.com/solarhell/go-deepl
```

## 快速开始

```go
client := deepl.New("your-auth-key")

translated, sourceLang, err := client.Translate(
    context.TODO(),
    "Hello, world.",
    deepl.Chinese,
)
```

如果 API key 为空，会自动读取 `DEEPL_AUTH_KEY` 环境变量。

## 翻译

```go
translated, sourceLang, err := client.Translate(
    context.TODO(),
    "Hello, world.",
    deepl.Chinese,
    deepl.SourceLang(deepl.English),
    deepl.Formality(deepl.LessFormal),
    deepl.TagHandling(deepl.HTMLTagHandling),
)
```

## 批量翻译

```go
translations, err := client.TranslateMany(
    context.TODO(),
    []string{"Hello, world.", "Goodbye."},
    deepl.Chinese,
)
for _, t := range translations {
    log.Printf("[%s] %s", t.DetectedSourceLanguage, t.Text)
}
```

## 客户端配置

```go
client := deepl.New("your-auth-key",
    deepl.BaseURL(deepl.FreeV2),
    deepl.HTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

## 错误处理

```go
translated, _, err := client.Translate(ctx, text, deepl.Chinese)
if err != nil {
    if deeplError, ok := errors.AsType[deepl.Error](err); ok {
        switch {
        case deeplError.IsBadRequest():
            // 请求参数无效 (400)
        case deeplError.IsUnauthorized():
            // API key 无效 (403)
        case deeplError.IsRateLimit():
            // 触发限流 (429)
        case deeplError.IsQuotaExceeded():
            // 字符额度已用尽 (456)
        case deeplError.IsServiceUnavailable():
            // 服务不可用 (503)
        }
    }
}
```

## 术语表管理

```go
// 创建术语表
glossary, err := client.CreateGlossary(ctx, "My Glossary",
    deepl.English, deepl.German,
    []deepl.GlossaryEntry{
        {Source: "hello", Target: "hallo"},
        {Source: "world", Target: "Welt"},
    },
)

// 列出所有术语表
glossaries, err := client.ListGlossaries(ctx)

// 删除术语表
err := client.DeleteGlossary(ctx, glossaryID)
```

## 测试

```bash
go test -short -race ./...
```

运行集成测试（**会产生 API 调用费用**）：

```bash
DEEPL_AUTH_KEY=YOUR_AUTH_KEY go test -race ./...
```

## 链接

- [DeepL API 文档](https://developers.deepl.com/docs)

## 许可证

[MIT](./LICENSE)
