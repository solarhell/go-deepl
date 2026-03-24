package deepl

import "time"

type translateResponse struct {
	Translations []Translation `json:"translations"`
}

// Translation is a translation result from deepl.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
	// BilledCharacters has the value only if ShowBilledChars(true) option was set.
	BilledCharacters int `json:"billed_characters"`
	// ModelTypeUsed has the value only if ProV2 api was used.
	ModelTypeUsed string `json:"model_type_used"`
}

// Glossary as per
// https://www.deepl.com/docs-api/managing-glossaries/creating-a-glossary/
type Glossary struct {
	GlossaryID   string    `json:"glossary_id"`
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	SourceLang   string    `json:"source_lang"`
	TargetLang   string    `json:"target_lang"`
	CreationTime time.Time `json:"creation_time"`
	EntryCount   int       `json:"entry_count"`
}

// A GlossaryEntry represents a single source→target entry in a glossary. This
// is serialized to/from tab-separated values for DeepL.
type GlossaryEntry struct {
	Source string
	Target string
}

// DocumentUploadResponse represents the response from the file upload API.
type DocumentUploadResponse struct {
	DocumentID  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

// DocumentStatusResponse represents the response from the file status API.
type DocumentStatusResponse struct {
	DocumentID       string `json:"document_id"`
	Status           string `json:"status"`
	SecondsRemaining int    `json:"seconds_remaining,omitempty"`
	BilledCharacters int    `json:"billed_characters,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

type UsageAndQuotaResponse struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}
