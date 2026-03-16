package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// DuckDuckGoBaseURL is exported so tests can override the endpoint.
var DuckDuckGoBaseURL = "https://html.duckduckgo.com/html/?q=%s"

// WebSearchHandler executes a DuckDuckGo HTML search and returns formatted results.
func WebSearchHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	query, err := requireStringArg(args, "query")
	if err != nil {
		return "", err
	}

	count := 5
	if c, ok := args["count"].(float64); ok {
		if c > 0 && c <= 10 {
			count = int(c)
		}
	}

	searchURL := fmt.Sprintf(DuckDuckGoBaseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return fmt.Sprintf("request creation error: %v", err), nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("request failed: %v", err), nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("HTTP Error %d", resp.StatusCode), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed reading response: %v", err), nil
	}

	htmlStr := string(bodyBytes)

	reLink := regexp.MustCompile(`<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>([\s\S]*?)</a>`)
	matches := reLink.FindAllStringSubmatch(htmlStr, count+5)

	if len(matches) == 0 {
		return "No results found.", nil
	}

	reSnippet := regexp.MustCompile(`<a class="result__snippet[^"]*".*?>([\s\S]*?)</a>`)
	snippetMatches := reSnippet.FindAllStringSubmatch(htmlStr, count+5)

	maxItems := count
	if len(matches) < count {
		maxItems = len(matches)
	}

	var results []string
	results = append(results, fmt.Sprintf("Results for: %s", query))

	for i := 0; i < maxItems; i++ {
		urlStr := matches[i][1]
		title := stripTags(matches[i][2])

		results = append(results, fmt.Sprintf("%d. %s\n   %s", i+1, title, decodeDuckDuckGoURL(urlStr)))

		if i < len(snippetMatches) {
			snippet := stripTags(snippetMatches[i][1])
			if snippet != "" {
				results = append(results, fmt.Sprintf("   %s", snippet))
			}
		}
	}

	return strings.Join(results, "\n"), nil
}

func stripTags(content string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(content, ""))
}

func decodeDuckDuckGoURL(raw string) string {
	if strings.Contains(raw, "uddg=") {
		if decoded, err := url.QueryUnescape(raw); err == nil {
			if idx := strings.Index(decoded, "uddg="); idx != -1 {
				return decoded[idx+5:]
			}
		}
	}
	return raw
}
