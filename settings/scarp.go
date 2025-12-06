package settings

import (
    "fmt"
    "math/rand"
    "net/http"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func ScrapWallpapers(query string) []string {
    imagesData := make([]string, 0)
    url := "https://wallpapers.com/search/" + query
    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching the URL:", err)
        return imagesData
    }
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        fmt.Println("Error loading the HTML:", err)
        return imagesData
    }

    totalPages := 1
    pageCounter := doc.Find(".page-counter.mobi")
    if len(pageCounter.Nodes) > 0 {
        totalPages, _ = strconv.Atoi(strings.Fields(pageCounter.Text())[len(strings.Fields(pageCounter.Text()))-1])
    }

    page := rand.Intn(totalPages) + 1
    pageURL := url + "?p=" + strconv.Itoa(page)

    response, err = http.Get(pageURL)
    if err != nil {
        fmt.Println("Error fetching the page URL:", err)
        return imagesData
    }
    defer response.Body.Close()

    doc, err = goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        fmt.Println("Error loading the HTML:", err)
        return imagesData
    }

    doc.Find("li.content-card").Each(func(i int, s *goquery.Selection) {
        imgTag := s.Find("img")
        imgURL := imgTag.AttrOr("data-src", "")
        if imgURL != "" {
            imageURL := strings.Join(strings.Split("https://wallpapers.com/", "/")[:len(strings.Split(pageURL, "/"))-1], "/") + imgURL
            imagesData = append(imagesData, imageURL)
        }
    })

    return imagesData
}
