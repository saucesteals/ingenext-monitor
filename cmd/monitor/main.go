package main

import (
	"log"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/saucesteals/ingenext-monitor"
)

func main() {
	hook, err := webhook.NewWithURL(os.Getenv("WEBHOOK_URL"))
	if err != nil {
		log.Panic(err)
	}

	versions, err := ingenext.GetVersions()
	if err != nil {
		log.Panic(err)
	}

	cache := versions

	// just as a confirmation for first attempt
	for v := range cache {
		cache[v] = cache[v][1:]
	}

	for {
		versions, err := ingenext.GetVersions()
		if err != nil {
			log.Printf("failed to get latest versions: %s", err)
			time.Sleep(time.Second * 10) // retry delay
			continue
		}

		for title, versions := range versions {
			added, removed := ingenext.VersionsDiff(cache[title], versions)

			if len(added) == 0 && len(removed) == 0 {
				continue
			}

			_, err = hook.CreateEmbeds([]discord.Embed{ingenext.CreateEmbed(title, added, removed)})
			if err != nil {
				log.Printf("failed to send webhook for %s: %s", title, err)
				continue
			}
		}

		cache = versions

		time.Sleep(time.Minute * 5)
	}

}
