## Intro

This is a console app that aggregates RSS feeds from the Internet. It has multi user capacity and can print aggregated feeds into the standard output.

## Prerequisites

- PostgreSQL
- Go 1.24+

## Installation

`go install gator`

## Config setup

Create `.gatorconfig.json` at your home directory.

### Config structure

```json
{
  "db_url": "postgres://YOUR_USE_NAME:@localhost:5432/gator?sslmode=disable",
  "current_user_name": "EXAMPLE_USER_NAME"
}
```

## Commands

### register

Registers a new user and sets the registered user as the current one.

Synatax: `gator register "user name"`


### users 

Lists all registered users.

Syntax: `gator users`

### login
Sets previously registered use as the current one.

Synatax: `gator login "User Name"`


### addfeed

Adds the named feed and follows it for the current user.

Synatax: `gator addfeed "Feed Name" "https:\\example.com/rss"`


### follow

Follows the named feed for the current user.

Synatax: `gator follow "Feed Name"`


### unfollow

Unfollows the named feed for the current user.

Synatax: `gator follow "Feed Name"`


### following

Lists all the feed that the current user follows.

Syntax: `gator following`

### browse

Prints the posts from the named feed to standard output. By default prints 2 posts, If the number of posts is set – prints the set number of posts.

Syntax:
* `gator browse "Feed Name"`
* `browse "Feed Name" 3` 

### agg

Starts the aggre`gator` service of the previously added feeds with a set time interval.

Sytax: 
* `gator agg 1s`
* `gator agg 1m`
* `gator agg 1h`
* `gator agg 1d`


### reset

Removes all the users, their feeds and the posts.

