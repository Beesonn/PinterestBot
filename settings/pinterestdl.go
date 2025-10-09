package settings 

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DownloadPinterestImage(URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	
	var imageURL string
	doc.Find("img.h-image-fit, img.h-unsplash-img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if src, exists := s.Attr("src"); exists {
			imageURL = src
			return false // break
		}
		return true
	})

	if imageURL == "" {
		doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if src, exists := s.Attr("src"); exists && strings.HasPrefix(src, "https://i.pinimg.com/") {
				imageURL = src
				return false // break
			}
			return true
		})
	}

	if imageURL == "" {
		return "", fmt.Errorf("no image found")
	}
	
	parts := strings.Split(imageURL, "x/")
	if len(parts) < 2 {
		return imageURL, nil 
	}
	secondPart := parts[1]
	newURL := "https://i.pinimg.com/1200x/" + secondPart

	return newURL, nil
}
