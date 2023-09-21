# How to setup and retrieve Github Token from 1Pass

Pre-requisite
- see [Get started with 1Password CLI](https://developer.1password.com/docs/cli/get-started/#step-1-install-1password-cli)
- 1Pass CLI client is installed (`brew install 1password-cli`)
- 1Pass CLI client is allowed to access database


## create token in GitHub
- see notes in Obsidian
- add it in 1Pass

## Load token in environment var

```sh
export GITHUB_TOKEN=$(op read "op://Dev - hacking/Jenkins_stats_PAT token/token")
```