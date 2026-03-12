# go-deepl

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/solarhell/go-deepl?style=for-the-badge)](https://goreportcard.com/report/github.com/solarhell/go-deepl)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/solarhell/go-deepl)
[![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/solarhell/go-deepl?utm_source=oss&utm_medium=github&utm_campaign=solarhell%2Fgo-deepl&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)](https://coderabbit.ai)

English | [中文](README_zh.md)

A thin, zero-dependency Go client for the [DeepL API](https://developers.deepl.com/docs). Built for Go 1.26+.

## Installation

```bash
go get github.com/solarhell/go-deepl
```

## Quick Start

```go
client := deepl.New("your-auth-key")

translated, sourceLang, err := client.Translate(
    context.TODO(),
    "Hello, world.",
    deepl.Chinese,
)
```

If the API key is empty, it reads from the `DEEPL_AUTH_KEY` environment variable.

## Translate

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

## Translate Multiple Texts

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

## Client Options

```go
client := deepl.New("your-auth-key",
    deepl.BaseURL(deepl.FreeV2),
    deepl.HTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

## Error Handling

```go
translated, _, err := client.Translate(ctx, text, deepl.Chinese)
if err != nil {
    if deeplError, ok := errors.AsType[deepl.Error](err); ok {
        switch {
        case deeplError.IsBadRequest():
            // invalid parameters (400)
        case deeplError.IsUnauthorized():
            // invalid API key (403)
        case deeplError.IsRateLimit():
            // rate limited (429)
        case deeplError.IsQuotaExceeded():
            // character limit reached (456)
        case deeplError.IsServiceUnavailable():
            // service unavailable (503)
        }
    }
}
```

## Glossary Management

```go
// Create a glossary
glossary, err := client.CreateGlossary(ctx, "My Glossary",
    deepl.English, deepl.German,
    []deepl.GlossaryEntry{
        {Source: "hello", Target: "hallo"},
        {Source: "world", Target: "Welt"},
    },
)

// List all glossaries
glossaries, err := client.ListGlossaries(ctx)

// Delete a glossary
err := client.DeleteGlossary(ctx, glossaryID)
```

## Testing

```bash
go test -short -race ./...
```

Run integration tests against the real DeepL API (**this will be billed**):

```bash
DEEPL_AUTH_KEY=YOUR_AUTH_KEY go test -race ./...
```

## Links

- [DeepL API Documentation](https://developers.deepl.com/docs)

## License

[MIT](./LICENSE)
