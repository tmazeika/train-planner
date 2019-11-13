# Train Planner
Scrapes [amtrak.com](https://www.amtrak.com) for train information and plans trips.

### Usage
In your cloned repository, run: `go run main.go [command]`.
##### Commands:
- **list \<from station\> \<to station\> \[date\]** Lists all trains between two stations for a given date (format '1/13/19') or today if empty.
- **plan \<from station\> \<to station\> \[date\]** Show journey plans between two stations for a given date (format '1/13/19') or today if empty. Supports `BBY` and `PHL` destinations/origins. Useful for Amtrak passriders wishing to book non-coach tickets as early as possible from the origin station of the "\<from station\>" argument.
- **save \<from station\> \<to station\> \[date\]** Saves the HTML of the fetched page for a trip between two stations for a given date (format '1/13/19') or today if empty. Saved in `trains.html`.
- **help \[command\]** Shows a list of commands or help for one command.

### Examples
- Find trains from Philadelphia's 30<sup>th</sup> Street Station to Boston's Back Bay Station on October 12, 2020: `go run main.go list phl bby 10/12/20`

### Notes
Web scraping is inherently unreliable. Breaking changes could happen at any time. Even a poor network connection could result in no results.

Caching will happen for "list" and "plan" commands. The `.trains.cache` file in the current working directory will be created/used/updated.

###### <sub>Neither I nor this project are affiliated, associated, authorized, endorsed by, or in any way officially connected with Amtrak, or any of its subsidiaries or its affiliates.</sub>
