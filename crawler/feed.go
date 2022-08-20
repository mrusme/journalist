package crawler

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"errors"
)

func (c *Crawler) GetFeedLink() (string, string, error) {
  if c.source == nil {
    return "", "", errors.New("No source available!")
  }

  if err := c.FromAuto(); err != nil {
    return "", "", nil
  }

  if c.contentType == "" {
    if err := c.Detect(); err != nil {
      return "", "", err
    }

    if c.contentType == "" {
      return "", "", errors.New("Could not detect content type!")
    }
  }

  if strings.Contains(c.contentType, "text/xml") {
    return "", c.sourceLocation, nil
  } else if strings.Contains(c.contentType, "text/html") {
    return c.GetFeedLinkFromHTML()
  }

  return "", "", errors.New("No feed link found")
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

