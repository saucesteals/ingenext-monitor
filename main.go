package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

var (
	urlUpdates    = "https://ingenext.ca/pages/safe-tesla-updates-for-boost50-and-bonus-module"
	cache         VersionHistory
	hook          webhook.Client
	delay         = time.Minute * 5
	isCron        = os.Getenv("CRON") == "1"
	diskCachePath = os.Getenv("VERSIONS_CACHE_PATH")
)

func sendWebhook(title string, added, removed []string) error {
	description := "New Changes:\n```diff"
	for _, ver := range added {
		description += fmt.Sprintf("\n+ %s", ver)
	}
	for _, ver := range removed {
		description += fmt.Sprintf("\n- %s", ver)
	}
	description += "```"

	now := time.Now()

	_, err := hook.CreateEmbeds([]discord.Embed{
		{
			Title:       title,
			URL:         urlUpdates,
			Description: description,
			Timestamp:   &now,
			Color:       0xf84248,
			Footer: &discord.EmbedFooter{
				Text:    "Ingenext Monitor",
				IconURL: "https://i.imgur.com/nJVc4DZ.jpg",
			},
		},
	})
	return err
}

func getDiskCache() (VersionHistory, error) {
	f, err := os.Open(diskCachePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var versions VersionHistory
	if err := json.NewDecoder(f).Decode(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func writeDiskCache(versions VersionHistory) error {
	data, err := json.Marshal(versions)
	if err != nil {
		return err
	}

	return os.WriteFile(diskCachePath, data, 0644)
}

func main() {
	var err error
	hook, err = webhook.NewWithURL(os.Getenv("WEBHOOK_URL"))
	if err != nil {
		log.Panic(err)
	}

	if diskCachePath != "" {
		versions, err := getDiskCache()
		if err != nil {
			log.Panic(err)
		}

		cache = versions
	} else {
		versions, err := getVersions()
		if err != nil {
			log.Panic(err)
		}

		cache = versions

		// just as a confirmation for first attempt
		for v := range cache {
			cache[v] = cache[v][1:]
		}
	}

	for {
		versions, err := getVersions()
		if err != nil {
			if isCron {
				log.Panic(err)
			}
			log.Printf("failed to get latest versions: %s", err)
			time.Sleep(time.Second * 10) // retry delay
			continue
		}

		for title, versions := range versions {
			added, removed := findUpdates(cache[title], versions)

			if len(added) == 0 && len(removed) == 0 {
				continue
			}

			if err := sendWebhook(title, added, removed); err != nil {
				if isCron {
					log.Panic(err)
				}
				log.Printf("failed to send webhook for %s: %s", title, err)
				continue
			}
		}

		if isCron {
			err := writeDiskCache(versions)
			if err != nil {
				log.Panic(err)
			}
			return
		}

		cache = versions

		time.Sleep(delay)
	}

}
