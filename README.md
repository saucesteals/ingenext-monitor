# Ingenext Safe Version Monitor

### Get alerted of all updates to the [Ingenext Safe Versions](https://ingenext.ca/pages/safe-tesla-updates-for-boost50-and-bonus-module)

## Running with Github Actions

- [Fork](https://github.com/saucesteals/ingenext-monitor/fork) this repo
- Go to the forked repo's actions settings (Setting > Secrets and variables > Actions)
- Create a secret `WEBHOOK_URL` and set it to your webhook's url
- Go to the **Actions** tab on the forked repo and enable workflows
- (OPTIONAL) Adjust the cron job schedule to your liking

## Manual usage

- Run (a built binary or with `go run ./cmd/monitor`) with the `WEBHOOK_URL` environment variable set to a discord webhook's url

## **License**

Distributed under the MIT License. See `LICENSE` for more information.
