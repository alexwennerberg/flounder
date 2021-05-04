// Parses Gemfeed according to the companion spec:
// gemini://gemini.circumlunar.space/docs/companion/subscription.gmi
package main

import (
	"bufio"
	"bytes"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
)

type Gemfeed struct {
	Title   string
	Creator string
	Url     *url.URL
	Entries []FeedEntry
}

func (gf *Gemfeed) toAtomFeed() string {
	feed := feeds.Feed{
		Title:  gf.Title,
		Author: &feeds.Author{Name: gf.Creator},
		Link:   &feeds.Link{Href: gf.Url.String()},
	}
	feed.Items = []*feeds.Item{}
	for _, fe := range gf.Entries {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   fe.Title,
			Link:    &feeds.Link{Href: fe.Url.String()}, // Rel=alternate?
			Created: fe.Date,                            // Updated not created?
			Content: fe.Content,
		})
	}
	res, _ := feed.ToAtom()
	return res
}

type FeedEntry struct {
	Title      string
	Url        *url.URL
	Date       time.Time
	DateString string
	Feed       *Gemfeed
	File       string // TODO refactor
	Content    string
}

func urlFromPath(fullPath string) url.URL {
	creator := getCreator(fullPath)
	baseUrl := url.URL{}
	baseUrl.Host = creator + "." + c.Host
	baseUrl.Path = getLocalPath(fullPath)
	return baseUrl
}

// Non-standard extension
// Requires yyyy-mm-dd formatted files
func generateFeedFromUser(user string) *Gemfeed {
	gemlogFolderPath := path.Join(c.FilesDirectory, user, GemlogFolder)
	// NOTE: assumes sanitized input
	u := urlFromPath(gemlogFolderPath)
	feed := Gemfeed{
		Title:   strings.Title(user) + "'s Gemlog",
		Creator: user,
		Url:     &u,
	}
	err := filepath.Walk(gemlogFolderPath, func(thepath string, info os.FileInfo, err error) error {
		base := path.Base(thepath)
		if len(base) >= 10 {
			entry := FeedEntry{}
			date, err := time.Parse("2006-01-02", base[:10])
			if err != nil {
				return nil
			}
			entry.Date = date
			entry.DateString = base[:10]
			entry.Feed = &feed
			content, err := ioutil.ReadFile(thepath)
			if err != nil {
				return nil
			}
			parse, _ := gemini.ParseText(bytes.NewReader(content))
			htmlDoc := textToHTML(nil, parse)
			if htmlDoc.Title != "" {
				// look for a title in headings
				entry.Title = htmlDoc.Title
			} else {
				// look for a title in lines
				f, err := os.Open(thepath)
				if err != nil {
					return nil
				}
				defer f.Close()
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					// skip blank lines
					if scanner.Text() == "" {
						continue
					}
					line := scanner.Text()
						var title string
						if len(line) > 50 {
							title = line[:50]
						} else {
							title = line
						}
						entry.Title = "[" + title + "...]"
					break
				}
			}
			entry.Content = htmlDoc.Content
			entry.File = getLocalPath(thepath)
			u := urlFromPath(thepath)
			entry.Url = &u
			feed.Entries = append(feed.Entries, entry)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	// Reverse chronological sort
	sort.Slice(feed.Entries, func(i, j int) bool {
		return feed.Entries[i].Date.After(feed.Entries[j].Date)
	})
	return &feed
}
