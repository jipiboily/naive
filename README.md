# PaaS from scratch

Some playground for my PaaS from scratch blog post or serie, part of an hypothetical blog serie, we'll see if I end up finishing/doing it!

If this end up being a thing, you'll see more on my blog, subscribe to my RSS feed: http://jipiboily.com :)

Note that this will most probably be a single node PaaS.

## The flow

My current idea for a basic version of this...

### Deploy process

- git push
- requiring a Dockerfile
- docker build ^
- check if app is started and up (paas.yaml, with a URL to check or something, and what to check…like “SERVER: OK”
- switch the previous one with new one (possibly change name and/or link? No idea, yet!)
- profit?

## Some cool ideas I will probably not do:

- client tool to
  - create new app
  - set env vars
  - restart
  - run one off containers
  - create databases and attach them to an app
  - scale ?
- coreos?