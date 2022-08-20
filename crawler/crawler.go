package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"golang.org/x/net/html"

	// "golang.org/x/net/html/charset"
	"errors"

	"golang.org/x/net/publicsuffix"

	scraper "github.com/tinoquang/go-cloudflare-scraper"
)

var DEFAULT_USER_AGENT string =
  "Mozilla/5.0 AppleWebKit/537.36 " +
  "(KHTML, like Gecko; compatible; " +
  "Googlebot/2.1; +http://www.google.com/bot.html)"

type Crawler struct {
  source         io.ReadCloser
  sourceLocation string
  contentType    string
}

func New() (*Crawler) {
  crawler := new(Crawler)

  crawler.source = nil
  crawler.Reset()
  return crawler
}

func (c *Crawler) Reset() {
  if c.source != nil {
    c.source.Close()
    c.source = nil
  }
  c.sourceLocation = ""
  c.contentType = ""
}

func (c *Crawler) FromHTTP(
  rawUrl *string,
  userAgent *string,
) (error) {
  if userAgent == nil {
    userAgent = &DEFAULT_USER_AGENT
  }

  jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
  if err != nil {
    return err
  }

  scraper, err := scraper.NewTransport(http.DefaultTransport)
  client := &http.Client{
    Jar: jar,
    Transport: scraper,
  }

  req, err := http.NewRequest("GET", *rawUrl, nil)
  if err != nil {
    return err
  }

  req.Header.Set("User-Agent",
    *userAgent)
  req.Header.Set("Accept",
    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif," +
      "image/webp,*/*;q=0.8")
  req.Header.Set("Accept-Language",
    "en-US,en;q=0.5")
  req.Header.Set("DNT",
    "1")

  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  // defer resp.Body.Close()

  c.Reset()
  c.source = resp.Body
  c.sourceLocation = *rawUrl
  return nil
}

func (c *Crawler) FromFile(rawUrl *string) (error) {
  file, err := os.Open(*rawUrl)
  if err != nil {
    return err
  }

  c.Reset()
  c.source = file
  return nil
}

func (c *Crawler) Detect() (error) {
  buf := make([]byte, 512)
  _, err := c.source.Read(buf)
  if err != nil {
    return err
  }

  c.contentType = http.DetectContentType(buf)
  return nil
}

func (c *Crawler) GetFeedLink() (string, error) {
  if c.source == nil {
    return "", errors.New("No source available!")
  }

  if c.contentType == "" {
    if err := c.Detect(); err != nil {
      return "", err
    }

    if c.contentType == "" {
      return "", errors.New("Could not detect content type!")
    }
  }

  return "", nil
}

func (c *Crawler) GetContentType() string {
  return c.contentType
}

func (c *Crawler) GetFeedLinkFromHTML() (string, string, error) {
  doc, err := html.Parse(c.source)
  if err != nil {
    return "", "", err
  }

  var f func(*html.Node) (bool, string, string)
  f = func(n *html.Node) (bool, string, string) {
    if n.Type == html.ElementNode && n.Data == "link" {
      var feedType *string = nil
      var feedHref *string = nil

      for i := 0; i < len(n.Attr); i++ {
        attr := n.Attr[i]
        if attr.Key == "type" {
          if strings.Contains(attr.Val, "rss") || strings.Contains(attr.Val, "atom") {
            feedType = &attr.Val
          }
        } else if attr.Key == "href" {
          feedHref = &attr.Val
        }
      }

      if feedType != nil && feedHref != nil {
        fmt.Printf("%v\n", *feedType)
        fmt.Printf("%v\n", *feedHref)
        return true, *feedType, *feedHref
      }

      return false, "", ""
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      fF, fT, fH := f(c)
      if fF == true {
        return fF, fT, fH
      }
    }
    return false, "", ""
  }

  found, feedType, feedHref := f(doc)
  if found == true {
    if strings.HasPrefix(feedHref, "./") {
      feedHref = fmt.Sprintf(
        "%s/%s",
        strings.TrimRight(c.sourceLocation, "/"),
        strings.TrimLeft(feedHref, "./"),
      )
    } else if strings.HasPrefix(feedHref, "/") {
      feedHref = fmt.Sprintf(
        "%s/%s",
        strings.TrimRight(c.sourceLocation, "/"),
        strings.TrimLeft(feedHref, "/"),
      )
    }
    return feedType, feedHref, nil
  }

  return "", "", errors.New("No feed URL found!")
}

