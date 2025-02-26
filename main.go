package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	flag.Parse()

	if *token == "" || *user == "" {
		fmt.Println("Error: Token & User ID required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	os.MkdirAll(*out, 0755)
	fmt.Printf("%d posts processed to %s\n", getAndProcessPosts(*user, *token, *out), *out)
}

func getAndProcessPosts(user, token, outDir string) int {
	count, url := 0, fmt.Sprintf("https://graph.instagram.com/v22.0/%s/media?fields=id,permalink,timestamp&access_token=%s", user, token)
	re := regexp.MustCompile(`instagram\.com/p/([^/]+)`)

	for url != "" {
		// API request and parsing
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			break
		}

		var result struct {
			Data []struct {
				ID, Permalink, Timestamp string
			} `json:"data"`
			Paging struct {
				Next string `json:"next,omitempty"`
			} `json:"paging"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			break
		}
		resp.Body.Close()

		// Process posts
		for _, p := range result.Data {
			id := strings.Split(p.ID, "_")[0] // First ID component
			sc := id                          // Shortcode (Fallback = ID)
			if m := re.FindStringSubmatch(p.Permalink); len(m) > 1 {
				sc = m[1]
			}

			t, _ := time.Parse("2006-01-02T15:04:05-0700", p.Timestamp)
			if t.IsZero() {
				t = time.Now()
			}

			fn := filepath.Join(outDir, fmt.Sprintf("%s-%s.md", t.Format("2006-01-02"), id))
			if err := os.WriteFile(fn, []byte(fmt.Sprintf(`---
date: %s
draft: false
tags:
  - instagram
---

{{< instagram %s >}}`, t.Format("2006-01-02T15:04:05-07:00"), sc)), 0644); err == nil {
				count++
			}
		}

		url = result.Paging.Next
		if url != "" {
			time.Sleep(time.Second)
		}
	}
	return count
}
