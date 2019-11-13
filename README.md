# Train Planner
Scrapes [amtrak.com](https://www.amtrak.com) for train information and plans trips.

### Usage
In your cloned repository, run: `go run main.go [command]`.
##### Commands:
- **list \<from station\> \<to station\> \[date\]** Lists all trains between two stations for a given date (format '1/13/19') or today if empty.
- **plan \<from station\> \<to station\> \[date\]** Show journey plans between two stations for a given date (format '1/13/19') or today if empty. Supports `BBY` and `PHL` destinations/origins.
- **save \<from station\> \<to station\> \[date\]** Saves the HTML of the fetched page for a trip between two stations for a given date (format '1/13/19') or today if empty. Saved in `trains.html`.
