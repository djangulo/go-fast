# go-fast

This is the learn-go-with-tests application, fully deployed through CI/CD

## Test

    $ docker-compose -f local.yml run --rm app go test -v

**IMPORTANT NOTE**

The production build does not include a go installation, only the compiled app binary; it will be impossible to run unit tests on the production build

## Deploy

#### To serve locally

Build the app:

    ~$ docker-compose -f local.yml build

then:

    ~$ docker-compose -f local.yml up

You can set the environment variable `COMPOSE_FILE` to avoid typing it every time

    ~$ export COMPOSE_FILE=local.yml
    ~$ docker-compose up --build

(`up` with `--build` flag builds [if necessary] and serves in one command)

### Deployment Checklist:

- Target machine has `docker` (>=18.09) && `docker-compose` (>=1.23.1) installed
  - The version check allows for the `-H` param
  - See installation instructions for [docker-ce](https://docs.docker.com/install/linux/docker-ce/ubuntu/) and [docker-compose](https://docs.docker.com/compose/install/)
- Target user belongs to the `docker` group (so user can run `docker` without elevation)
  - Run `sudo usermod -aG docker $USER`

### Staging build

It runs a single instance, with a static `traefik` config. Requires no setup other than the commands below. This will be used for the E2E suite (not yet implemented).

Set the `DOCKER_HOST` environment variable

    $ export DOCKER_HOST=ssh://user@host
    $ docker-compose -f staging.yml up --detach --build

or run the `-H` parameter to `docker-compose`:

    $ docker-compose -H "ssh://user@host" -f staging.yml up --build --detach

### Production build

Load balanced, scalable build (not yet implemented)
