package crawler

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"strings"
	"time"
)

func Crawl() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("cannot install playwright: %v", err)
	}

	pw, err := playwright.Run()

	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	err = page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept-Language": "en-US,en;q=0.9",
	})
	if err != nil {
		log.Fatalf("could not set headers: %v", err)
		return
	}

	var toVisit = []string{"/"}
	var visited = make(map[string]struct{})

	for len(toVisit) > 0 {
		var p = toVisit[0]
		toVisit = toVisit[1:]

		if _, exists := visited[p]; !exists {
			if _, err := page.Goto(fmt.Sprintf("https://pcdiga.com%s", p)); err != nil {
				log.Fatalf("could not goto: %v", err)
			}

			visited[p] = struct{}{}

			pageLinks, err := fetchLinks(page)
			if err != nil {
				fmt.Printf("could not fetch page links for %s: %v\n", p, err)
			} else if len(pageLinks) > 0 {
				toVisit = append(toVisit, pageLinks...)
				fmt.Printf("added %d more links from %s, queue size: %d\n", len(pageLinks), p, len(toVisit))
			}

			time.Sleep(2 * time.Second)
		}

	}

	fmt.Printf("%s", toVisit)
}

func fetchProductInfo(page playwright.Page) {

}

func fetchLinks(page playwright.Page) ([]string, error) {
	entries, err := page.Locator("a").All()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		href, err := entry.GetAttribute("href")
		if err != nil {
			fmt.Println("No href, skip")
		}
		if strings.HasPrefix(href, "/") {
			result = append(result, href)
		}
	}

	return result, nil
}
