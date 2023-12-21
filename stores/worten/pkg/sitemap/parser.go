package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Adapted from https://github.com/Z-M-Huang/sitemap-parser/blob/master/parser.go

// Index sitemap index
type Index struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Elements []Element `xml:"sitemap"`
}

// Sitemap sitemap data
type Sitemap struct {
	XMLName  xml.Name  `xml:"urlset"`
	Elements []Element `xml:"url"`
}

// Element single sitemap element
type Element struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float32 `xml:"priority"`
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getSitemapFromBytes(body []byte) (*Sitemap, error) {
	ret := &Sitemap{}
	err := xml.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetIndex get sitemap index from URL
func GetIndex(url string) (*Index, error) {
	respBytes, err := get(url)
	if err != nil {
		return nil, err
	}
	ret := &Index{}
	err = xml.Unmarshal(respBytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetSitemap get sitemap from URL
func GetSitemap(url string) (*Sitemap, error) {
	respBytes, err := get(url)
	if err != nil {
		return nil, err
	}
	return getSitemapFromBytes(respBytes)
}

// GetSitemapGZ get sitemaps from .gz URL
func GetSitemapGZ(url string) (*Sitemap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return getSitemapFromBytes(bodyBytes)
}

// GetSitemaps loads all sitemaps from index
func (i *Index) GetSitemaps() ([]*Sitemap, error) {
	var sitemaps []*Sitemap

	for _, e := range i.Elements {
		if strings.HasSuffix(e.Loc, ".xml") {
			s, err := GetSitemap(e.Loc)
			if err != nil {
				return sitemaps, err
			}
			sitemaps = append(sitemaps, s)
		} else if strings.HasSuffix(e.Loc, ".gz") {
			//get .gz
			s, err := GetSitemapGZ(e.Loc)
			if err != nil {
				return sitemaps, err
			}
			sitemaps = append(sitemaps, s)
		} else {
			return sitemaps, fmt.Errorf("invalid sitemap loc: %s", e.Loc)
		}
	}
	return sitemaps, nil
}
