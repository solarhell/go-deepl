package deepl_test

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/solarhell/go-deepl"
)

func ExampleClient_Translate() {
	client := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))

	translated, sourceLang, err := client.Translate(context.TODO(), "Hello, world.", deepl.Chinese)
	if err != nil {
		if deeplError, ok := errors.AsType[deepl.Error](err); ok {
			log.Fatalf("deepl api error code %d: %s", deeplError.Code, deeplError.Error())
		}
		log.Fatal(err)
	}

	log.Printf("source language: %s", sourceLang)
	log.Println(translated)
}

func ExampleClient_Translate_withOptions() {
	client := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))

	translated, sourceLang, err := client.Translate(
		context.TODO(),
		"Hello, world.",
		deepl.Chinese,
		deepl.SourceLang(deepl.English),
		deepl.SplitSentences(deepl.SplitNoNewlines),
		deepl.PreserveFormatting(true),
		deepl.Formality(deepl.LessFormal),
	)
	if err != nil {
		if deeplError, ok := errors.AsType[deepl.Error](err); ok {
			log.Fatalf("deepl api error code %d: %s", deeplError.Code, deeplError.Error())
		}
		log.Fatal(err)
	}

	log.Printf("source language: %s", sourceLang)
	log.Println(translated)
}

func ExampleClient_TranslateMany() {
	client := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))

	translations, err := client.TranslateMany(
		context.TODO(),
		[]string{
			"Hello, world.",
			"This is an example.",
			"Goodbye.",
		},
		deepl.Chinese,
	)
	if err != nil {
		if deeplError, ok := errors.AsType[deepl.Error](err); ok {
			log.Fatalf("deepl api error code %d: %s", deeplError.Code, deeplError.Error())
		}
		log.Fatal(err)
	}

	for _, translation := range translations {
		log.Printf("source language: %s", translation.DetectedSourceLanguage)
		log.Println(translation.Text)
		log.Println()
	}
}

func ExampleClient_TranslateMany_withOptions() {
	client := deepl.New(os.Getenv("DEEPL_AUTH_KEY"))

	translations, err := client.TranslateMany(
		context.TODO(),
		[]string{
			"Hello, world.",
			"This is an example.",
			"Goodbye.",
		},
		deepl.Chinese,
		deepl.SourceLang(deepl.English),
		deepl.SplitSentences(deepl.SplitNoNewlines),
		deepl.PreserveFormatting(true),
		deepl.Formality(deepl.LessFormal),
	)
	if err != nil {
		if deeplError, ok := errors.AsType[deepl.Error](err); ok {
			log.Fatalf("deepl api error code %d: %s", deeplError.Code, deeplError.Error())
		}
		log.Fatal(err)
	}

	for _, translation := range translations {
		log.Printf("source language: %s", translation.DetectedSourceLanguage)
		log.Println(translation.Text)
		log.Println()
	}
}
