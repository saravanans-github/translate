// Sample translate-quickstart translates "Hello, world!" into Russian.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"srt"

	// Imports the Google Cloud Translate client package.
	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

var _lang string

func main() {
	inFile := flag.String("in", "", "SRT input file")
	lang := flag.String("la", "en", "Language to translate to")
	flag.Parse()

	filename := *inFile
	_lang = *lang

	readElements := make(chan srt.Element)
	writeElements := make(chan srt.Element)
	go srt.Read(filename, readElements)

	go func() {
		for element := range readElements {
			writeElement := element
			for i, subtitle := range element.Subtitles {
				element.Subtitles[i] = translateText(subtitle)
			}
			writeElements <- writeElement
		}
		close(writeElements)
	}()

	srt.Write("out_"+filename, writeElements)
}

func translateText(text string) string {
	ctx := context.Background()

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the text to translate.
	// Sets the target language.
	target, err := language.Parse(_lang)
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	// Translates the text into Tamil.
	translations, err := client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	fmt.Printf("Text: %v\n", text)
	fmt.Printf("Translation: %v\n", translations[0].Text)

	return translations[0].Text
}
