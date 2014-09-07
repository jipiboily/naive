# MVP

## What's missing?
- nginx (could be on host for now)
- postgres (could be on host for now)
- database backups
- support for Dockerfile OR buildstep
- support for custom domains
- tests
- better structure
  - single file still (for now), but not all should belong to the Release struct

## Right now, I have:
- Dockerfile being built
- docker run
  - will fail on name collision for now
- nginx running

## Next steps
- figure out a way to push new version without name conflict
- nginx: see the Rails app
  - update that, on new deploy
- database
- env vars
  - be able to change them in running container + for future

## Internal storage

### What to store?
- release history
- environment variables
- name of current containers running
- use $PORT to start stuff inside the container

### How?

- flat files for now I would say?

## When creating a new project
- create directory structure
- create the nginx config
- add a symbolic link to it, from the sites-enabled directory