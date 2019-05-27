# Git Manager

## How to use it

### Build

`go run make.go`

### Use

Add a repository:

`gm add ~/repositories/gocv.io/x/gocv`

Then, check your configuration:

`gm config -p` // -p means pretty print   
		
Then trigger an update for the repositories

`gm update`

You can also delete a repository

`gm delete gocv`
 
## When to use it

I did it to scratch my itch, while working with many small repositories for the same project I had to keep updating them individually, that was tedious and easy to lose track.

## What it is

It's a simple tool to manage many repositories without the hassle of doing one by one and manually.

## What it's not

It's not a git wrapper. So, keep using your git-fu on your projects.
