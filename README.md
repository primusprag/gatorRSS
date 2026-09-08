# Welcome to GatorRSS by d4l4-33

### This is a boot.dev guided project for a RSS aggregator called gatorRSS.

#### Required programs:

- Postgres
- Go

#### To install:

```
go install gatorRSS github.com/d4l4-33/gatorRSS
```


#### Commands:

- Users:
register "username" - Creates a user that stores your feeds
login "username" - Selects user
users - Lists the created users

- Feeds:
addfeed "feed name" "url" - Creates a feed and connects it to the current user
feeds- Lists feeds

- Follow:
following: - Lists feeds that the current user is following/recieving posts from
follow "feed name" - Connects an existing feed to the current user. Automatically assigned when using "addfeed"
unfollow "feed name" - Removes the feed from the list of feeds that the current user recieves

- Posts:
agg "interval" - Gathers posts from the feed longest since gatherd that the current user is following after an "interval" delay.
Examples: agg 10s (every 10 seconds), agg 30m (every 30 minutes), agg 5h (every 5 hours).
Browse "*amount*" - Lists the latest posts. Amount is 5 by default if amount not assigned.

-Archiving:
archive "name" - Archives a feed and it's posts

restore "name" - Resores a feed but not posts, use agg to gather new posts
-Resetting:
reset - Emptys the users and feeds database
resetposts - Emptys the posts but keeps the users and feeds.
