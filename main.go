package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	urlFlag := flag.String("url", "", "URL to scrape")
	extractFlag := flag.String("extract", "links", "Elements to extract (e.g., 'links', 'headlines', 'all')")
	outputFileFlag := flag.String("output", "", "Output file path (optional)")

	flag.Parse()

	targetURL := *urlFlag
	extractValue := *extractFlag
	outputFilePath := *outputFileFlag

	if targetURL == "" {
		fmt.Fprintln(os.Stderr, "Error: The -url flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	if !(strings.HasPrefix(targetURL, "http://") || strings.HasPrefix(targetURL, "https://")) {
		fmt.Fprintf(os.Stderr, "Error: Invalid URL provided via -url flag: %s\n", targetURL)
		fmt.Fprintln(os.Stderr, "URL must start with 'http://' or 'https://'.")
		os.Exit(1)
	}

	baseURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse base URL '%s': %v\n", targetURL, err)
		os.Exit(1)
	}
	fmt.Printf("Base URL successfully parsed: %s\n", baseURL.String())
	fmt.Printf("Extraction type set to: %s\n", extractValue)

	var file *os.File
	var fileErr error
	var outputWriter io.Writer = os.Stdout

	if outputFilePath != "" {
		fmt.Printf("Attempting to create/open output file: %s\n", outputFilePath)
		file, fileErr = os.Create(outputFilePath)
		if fileErr != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not create output file '%s': %v\n", outputFilePath, fileErr)
			os.Exit(1)
		}
		fmt.Printf("Successfully opened file %s for writing.\n", outputFilePath)
		defer file.Close()
		outputWriter = file
	} else {
		fmt.Println("Output will be printed to the console.")
	}

	fmt.Println("Attempting to fetch URL:", targetURL)
	resp, httpErr := http.Get(targetURL)
	if httpErr != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL %s: %v\n", targetURL, httpErr)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: Received non-200 status code %d for URL %s\n", resp.StatusCode, targetURL)
		os.Exit(1)
	}
	fmt.Println("HTTP request successful (Status 200 OK)! Parsing HTML body...")

	doc, parseDocErr := goquery.NewDocumentFromReader(resp.Body)
	if parseDocErr != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HTML for %s: %v\n", targetURL, parseDocErr)
		os.Exit(1)
	}
	fmt.Println("HTML parsed successfully! Ready to extract data...")

	var extractedLinks []string
	var extractedHeadlines []string

	if extractValue == "links" || extractValue == "all" {
		fmt.Println("Extracting links...")
		linkSelection := doc.Find("a")
		linkSelection.Each(func(index int, element *goquery.Selection) {
			hrefValue, hrefExists := element.Attr("href")
			if hrefExists {
				linkURL, parseLinkErr := url.Parse(hrefValue)
				if parseLinkErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: Skipping malformed link #%d: '%s' - Error: %v\n", index+1, hrefValue, parseLinkErr)
					return
				}
				absoluteURL := baseURL.ResolveReference(linkURL)
				absoluteURLString := absoluteURL.String()
				extractedLinks = append(extractedLinks, absoluteURLString)
			}
		})
		fmt.Printf("Successfully stored %d absolute link URL(s).\n", len(extractedLinks))
	}
	if extractValue == "headlines" || extractValue == "all" {
		fmt.Println("Extracting headlines...")
		headlineSelection := doc.Find("h1, h2, h3")
		headlineSelection.Each(func(index int, element *goquery.Selection) {
			headlineText := strings.TrimSpace(element.Text())
			if headlineText != "" {
				extractedHeadlines = append(extractedHeadlines, headlineText)
			}
		})
		fmt.Printf("Successfully stored %d non-empty headline(s).\n", len(extractedHeadlines))
	}

	if extractValue == "links" || extractValue == "all" {
		_, err := fmt.Fprintln(outputWriter, "\n--- Links ---")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to write links header to output: %v\n", err)
		}
		if len(extractedLinks) > 0 {
			for _, link := range extractedLinks {
				_, err := fmt.Fprintln(outputWriter, link)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to write link '%s' to output: %v\n", link, err)
				}
			}
		} else {
			_, err := fmt.Fprintln(outputWriter, "No links found or extracted.")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to write 'no links' message to output: %v\n", err)
			}
		}
	}

	if extractValue == "headlines" || extractValue == "all" {
		_, err := fmt.Fprintln(outputWriter, "\n--- Headlines ---")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to write headlines header to output: %v\n", err)
		}
		if len(extractedHeadlines) > 0 {
			for _, headline := range extractedHeadlines {
				_, err := fmt.Fprintln(outputWriter, headline)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to write headline '%s' to output: %v\n", headline, err)
				}
			}
		} else {
			_, err := fmt.Fprintln(outputWriter, "No headlines found or extracted.")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to write 'no headlines' message to output: %v\n", err)
			}
		}
	}
	if extractValue != "links" && extractValue != "headlines" && extractValue != "all" {
		fmt.Fprintf(os.Stderr, "\nWarning: Invalid value '%s' provided for -extract flag. ", extractValue)
		fmt.Fprintln(os.Stderr, "Valid options are 'links', 'headlines', or 'all'. No data extracted/printed.")
	}
	fmt.Println("\nScraping process finished.")
}