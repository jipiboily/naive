# PaaS from scratch

Some playground for my PaaS from scratch blog post or serie, part of an hypothetical blog serie, we'll see if I end up finishing/doing it! It will not be really from scratch, but using some quite low-level-ish pieces not everyone knows about. A real from scratch here would be doable but would require more work than the time I have to spend on this.

If this end up being a thing, you'll see more on my blog, subscribe to my RSS feed: http://jipiboily.com :)

Note that this will most probably be a single node PaaS.

## The flow

My current idea for a basic version of this...

### Deploy process

- git push
- how to build it:
  - if there is a Dockerfile, use it?
  - else, use `buildstep` to build with heroku buildpacks
- docker build ^
- check if app is started and up (paas.yaml, with a URL to check or something, and what to check…like “SERVER: OK”
- switch the previous one with new one (possibly change name and/or link? No idea, yet!)
- profit?

### Possible directory structure

For each app, we need to have a directory with that kind of structure:

- project_name
  - current -> this is a symlink to the latest release
  - logs
    - history.log (including a timestamp, the new and old SHA1)
    - builds
      - $(SHA1).log -> using the SHA1 of the release as the name. This is the build log.
  - releases
    - one directory per release, with the SHA1 as the name, containing the app's source for that release (this might actually be useless, but let's do it for now.)

## Some cool ideas I will probably not do:

- client tool to
  - create/delete apps
  - set/unset env vars
  - restart/stop/start apps
  - **run one off containers** (that one is probably the most important, to be able to run Rails migrations as an example)
  - create databases and attach them to an app
  - scale ?
- coreos?

## Missing

- history of deploys and a way to roll back
