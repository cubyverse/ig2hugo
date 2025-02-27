package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Basic setup
	token := flag.String("token", "", "Token")
	user := flag.String("user", "", "User ID")
	out := flag.String("output", "content/instagram", "Output directory")
	debug := flag.Bool("debug", false, "Enable detailed debugging")
	flag.Parse()

	if *token == "" || *user == "" {
		fmt.Println("Error: Token & User ID required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	os.MkdirAll(*out, 0755)
	fmt.Printf("%d posts processed to %s\n", getAndProcessPosts(*user, *token, *out, *debug), *out)
}

func getAndProcessPosts(user, token, outDir string, debug bool) int {
	count, url := 0, fmt.Sprintf("https://graph.instagram.com/v22.0/%s/media?fields=id,permalink,timestamp,media_type,media_url&access_token=%s", user, token)
	re := regexp.MustCompile(`instagram\.com/p/([^/]+)`)
	reelRe := regexp.MustCompile(`instagram\.com/reel/([^/]+)`)

	for url != "" {
		// API request and parsing
		if debug {
			fmt.Printf("Fetching URL: %s\n", url)
		}

		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("Error fetching from API: %v, status: %d\n", err, resp.StatusCode)
			if resp != nil && resp.Body != nil {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: %s\n", string(body))
				resp.Body.Close()
			}
			break
		}

		var result struct {
			Data []struct {
				ID, Permalink, Timestamp string
				MediaType                string `json:"media_type,omitempty"`
				MediaURL                 string `json:"media_url,omitempty"`
			} `json:"data"`
			Paging struct {
				Next string `json:"next,omitempty"`
			} `json:"paging"`
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if debug {
			fmt.Printf("Raw API response: %s\n", string(body))
		}

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			break
		}

		// Process posts
		for _, p := range result.Data {
			id := strings.Split(p.ID, "_")[0] // First ID component
			sc := id                          // Shortcode (Fallback = ID)

			if debug {
				fmt.Printf("\nProcessing post:\n")
				fmt.Printf("  ID: %s\n", p.ID)
				fmt.Printf("  Permalink: %s\n", p.Permalink)
				fmt.Printf("  Timestamp: %s\n", p.Timestamp)
				fmt.Printf("  Media Type: %s\n", p.MediaType)
				fmt.Printf("  Media URL: %s\n", p.MediaURL)
			}

			if p.Permalink == "" {
				fmt.Printf("Warning: Post %s has no permalink, using ID as shortcode\n", id)
			} else {
				// Try to extract shortcode from permalink
				if m := re.FindStringSubmatch(p.Permalink); len(m) > 1 {
					sc = m[1]
					if debug {
						fmt.Printf("  Extracted shortcode from /p/ URL: %s\n", sc)
					}
				} else if m := reelRe.FindStringSubmatch(p.Permalink); len(m) > 1 {
					sc = m[1]
					if debug {
						fmt.Printf("  Extracted shortcode from /reel/ URL: %s\n", sc)
					}
				} else {
					fmt.Printf("Warning: Could not extract shortcode from permalink '%s' for post %s (type: %s), using ID\n",
						p.Permalink, id, p.MediaType)

					// Try to extract shortcode from the last part of the URL
					parts := strings.Split(strings.TrimRight(p.Permalink, "/"), "/")
					if len(parts) > 0 {
						possibleSc := parts[len(parts)-1]
						if possibleSc != "" && possibleSc != id {
							sc = possibleSc
							fmt.Printf("  Extracted possible shortcode from URL path: %s\n", sc)
						}
					}
				}
			}

			t, err := time.Parse("2006-01-02T15:04:05-0700", p.Timestamp)
			if err != nil {
				fmt.Printf("Error parsing timestamp '%s': %v, using current time\n", p.Timestamp, err)
				t = time.Now()
			}

			fn := filepath.Join(outDir, fmt.Sprintf("%s-%s.md", t.Format("2006-01-02"), id))
			content := fmt.Sprintf(`---
date: %s
draft: false
tags:
  - instagram
instagram_id: %s
instagram_permalink: %s
media_type: %s
shortcode: %s
---

{{< instagram %s >}}`, t.Format("2006-01-02T15:04:05-07:00"), id, p.Permalink, p.MediaType, sc, sc)

			if err := os.WriteFile(fn, []byte(content), 0644); err == nil {
				count++
				fmt.Printf("Created post %s with shortcode %s\n", fn, sc)
			} else {
				fmt.Printf("Error writing file %s: %v\n", fn, err)
			}
		}

		url = result.Paging.Next
		if url != "" {
			time.Sleep(time.Second)
		}
	}
	return count
}
