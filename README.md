# naive

Naive is a very naive and simple PaaS. It is only for a single server, not scallable, not secure, super hacky, but fun to play with.

## Why?

This is a playground for my "PaaS from scratch" blog post or serie, part of an hypothetical "from scratch" blog serie, we'll see if I end up finishing/doing it! It will not be really from scratch, but using some quite low-level-ish pieces not everyone knows about. A real from scratch here would be doable but would require more work than the time I have to spend on this.

If this end up being a thing, you'll see more on my blog, subscribe to my RSS feed: http://jipiboily.com :)

Note that this will most probably be a single node PaaS.

## How will that work?

This is not final, but my current idea for a basic version of this...

### Deploy process

- git push (app will be created if it doesn't exists)
- how to build it:
  - if there is a Dockerfile, use it!
  - else, use `buildstep` to build with heroku buildpacks
- docker build ^
- docker run the new version
- point nginx to the new container
- stop the old container
- profit?

## HOWTOS

### Install

### Create a new app

- ...
- To setup environment variables, create a file a `app.env` file at the root of the project's directory on the server: `sudo vi /paas/apps/my-project-name/app.env`. Each line will be used as is and should respect the "ENV_VAR=VALUE" format, no comment allowed or anything else.
- ...

## Some cool ideas I will probably not do:

- client tool to
  - create/delete apps
  - set/unset env vars
  - restart/stop/start apps
  - **run one off containers** (that one is probably the most important, to be able to run Rails migrations as an example)
  - create databases and attach them to an app
  - scale ?
- coreos
  - if this was a more serious/less naive project, that would be the way to go IMHO

## TODO & missing

- See TODO.md for now. (but this is not up to date, I am using a private Trello board right now. Let me know if you want to see it or help.)
