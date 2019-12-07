# Git Manager

	A little helper for your git repositories.

## Build and Install

`go run make.go`

## How to use it

- Add a repository

`gm add ~/repositories/gocv.io/x/gocv`

- Or clone remote repository and add to configuration ( the parent must exist )

`gm get https://github.com/upspin/upspin.git ~/repositories/upspin`

- Then, check your configuration

`gm config -p` // -p means pretty print

- Then trigger an update for the repositories

`gm update`

- You can also delete a repository

`gm delete gocv`

## When to use it

I did it to scratch my itch, while working with many small repositories for the same project I had to keep updating them individually, that was tedious and easy to lose track.

## What it is

It's a simple tool to manage many repositories without the hassle of doing one by one and manually.

## What it's not

It's not a git wrapper. So, keep using your git-fu on your projects.
