# contractus

![Logo](/assets/blackbackground.png#gh-dark-mode-only)
![Logo](/assets/whitebackgroun.png#gh-light-mode-only)

![image](assets/contractus.png)


# Contractus

Jamming around orders with API endpoints 🎸

## Environment Variables

**Required** environment variables

    CONTRACTUS_POSTGRES_URL
    CONTRACTUS_GOOGLE_CLIENT_ID
    CONTRACTUS_GOOGLE_CLIENT_SECRET
    CONTRACTUS_DOMAIN
    CONTRACTUS_JWT_SECRET_KEY

## Get Started

To start the service locally, you can type `make dev/start` and after that you can use the docker container IP to play around the routes, `make ip`

Request example:

    curl (make ip)

## Command line

All commands are synthesized in the Makefile `make help`, to start the development environment, just run:

    make dev/start
    make dev <- You will be able to run commands inside the container

After running `make dev`, it's possible to run the following commands inside the container:
    
- The integration tests: `make integration-test testcase=<>`
- The Unit tests: `make test testcase=<>`
- The lint: `make lint`

Or, it's also possible run from local:

- make `dev/integration-test testcase=<>`
- make `dev/test testcase=<>`

The testcase variable could be used to run a specific test

## Ship a new version
    `make image/publish`
    `heroku container:release web -a contractus`

## Logs

Production
    `heroku logs --tail -a contractus`

Local:
    `make dev/logs contractus`


## API documentation
[API Docs](api/docs/)

## Attention points 
 - For a while, the integration-tests just ran locally not in CI, this increased the time to ship code 🚀
 - We don't have a way to paginate transactions. 😔
 - To publish images and new releases, for now, the only way is using the command line, isn't automate by CI yet 😔
