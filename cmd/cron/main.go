package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/saucesteals/ingenext-monitor"
)

var (
	diskCachePath = os.Getenv("VERSIONS_CACHE_PATH")
)

func main() {
	hook, err := webhook.NewWithURL(os.Getenv("WEBHOOK_URL"))
	if err != nil {
		log.Panic(err)
	}

	cache, err := getDiskCache()
	if err != nil {
		log.Panic(err)
	}

	if cache == nil {
		versions, err := ingenext.GetVersions()
		if err != nil {
			log.Panic(err)
		}

		cache = versions

		// just as a confirmation for first attempt
		for v := range cache {
			cache[v] = cache[v][1:]
		}
	}

	versions, err := ingenext.GetVersions()
	if err != nil {
		log.Panic(err)
	}

	for title, versions := range versions {
		added, removed := ingenext.VersionsDiff(cache[title], versions)

		if len(added) == 0 && len(removed) == 0 {
			continue
		}

		_, err = hook.CreateEmbeds([]discord.Embed{ingenext.CreateEmbed(title, added, removed)})
		if err != nil {
			log.Panicf("failed to send webhook for %s: %s", title, err)
		}
	}

	err = writeDiskCache(versions)
	if err != nil {
		log.Panic(err)
	}
}

func getDiskCache() (ingenext.VersionHistory, error) {
	f, err := os.Open(diskCachePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	defer f.Close()

	var versions ingenext.VersionHistory
	if err := json.NewDecoder(f).Decode(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func writeDiskCache(versions ingenext.VersionHistory) error {
	f, err := os.OpenFile(diskCachePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewEncoder(f).Encode(versions)
}
