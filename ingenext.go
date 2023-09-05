package ingenext

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/disgoorg/disgo/discord"
)

var (
	urlUpdates = "https://ingenext.ca/pages/safe-tesla-updates-for-boost50-and-bonus-module"
)

type VersionHistory map[string][]string

func normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func GetVersions() (VersionHistory, error) {
	res, err := http.Get(urlUpdates)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	mainContent := doc.Find("#MainContent")

	versions := VersionHistory{}

	titles := mainContent.Find("p > strong")

	uls := mainContent.Find("ul").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return len(s.Nodes[0].Attr) == 0
	})

	titles.Each(func(i int, title *goquery.Selection) {
		if i == 0 {
			return
		}

		titleString := strings.TrimSpace(title.Text())

		ul := uls.Eq(i)

		lis := ul.Find("li")

		lis.Each(func(i int, li *goquery.Selection) {
			versions[titleString] = append(versions[titleString], normalize(li.Text()))
		})
	})

	return versions, nil
}

func VersionsDiff(cached, latest []string) (added, removed []string) {
	c := map[string]struct{}{}

	for _, ver := range cached {
		c[ver] = struct{}{}
	}

	for _, ver := range latest {
		_, exists := c[ver]
		if exists {
			delete(c, ver)
			continue
		}

		added = append(added, ver)
	}

	for ver := range c {
		removed = append(removed, ver)
	}

	return
}

func CreateEmbed(title string, added, removed []string) discord.Embed {
	description := "New Changes:\n```diff"
	for _, ver := range added {
		description += fmt.Sprintf("\n+ %s", ver)
	}
	for _, ver := range removed {
		description += fmt.Sprintf("\n- %s", ver)
	}
	description += "```"

	now := time.Now()

	return discord.Embed{
		Title:       title,
		URL:         urlUpdates,
		Description: description,
		Timestamp:   &now,
		Color:       0xf84248,
		Footer: &discord.EmbedFooter{
			Text:    "Ingenext Monitor",
			IconURL: "https://i.imgur.com/nJVc4DZ.jpg",
		},
	}
}
