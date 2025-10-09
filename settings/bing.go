package settings

import (
    "encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ImageObject struct {
	URL string `json:"url"` 
}

func SearchBing(keyword string, nbImages int) ([]ImageObject, error) {
	query := strings.ReplaceAll(keyword, " ", "+")
	url := fmt.Sprintf("https://www.bing.com/images/search?q=%s&safeSearch=strict", query)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var imagesObjects []ImageObject
	
	doc.Find(".iusc").Each(func(i int, s *goquery.Selection) {
		jsonContent, exists := s.Attr("m")
		if exists {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonContent), &result)
			if err == nil {
				if murl, ok := result["murl"].(string); ok {					
					imagesObjects = append(imagesObjects, ImageObject{URL: murl})
				}
			}
		}
	})
	
	if nbImages > 0 && nbImages < len(imagesObjects) {
		imagesObjects = imagesObjects[:nbImages]
	}

	return imagesObjects, nil
}
