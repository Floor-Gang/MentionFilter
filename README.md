# MentionFilter
Allows moderators to set regex filters for discord moderation. Eligible for direct removal or filter and possibly manual removal.

[Trello card - claimed by @mandjevant](https://trello.com/c/FtFfTVzh)

# Usage
The commands to use this service are as follows:

```
# Add a mention
.[syntax command] add <regex> <action> <description>

# Remove a mention
.[syntax command] remove <id>

# Display all mentions
.[syntax command] mentions

# Display a singular mention
.[syntax command] mention <id>

# Change what happens on mention
.[syntax command] change_action <id> <type (filter/remove)>

# Change regex of mention
.[syntax command] change_regex <id> <regex>

# Change description of mention
.[syntax command] change_description <id> <description>

# Display help message
.[syntax command] help
```

# Database layout

This will make the rules as dynamic as possible.

| mention_id | regex | action | description |
|------------|-------|--------|-------------|
| 1          | regex | action | description | 

# Setup
Download golang if you haven't already at https://golang.org/dl/ after that install the packages 

```
$ go get
$ go build 

```
