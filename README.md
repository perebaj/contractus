# contractus

![Logo](/assets/blackbackground.png#gh-dark-mode-only)
![Logo](/assets/whitebackgroun.png#gh-light-mode-only)

![image](assets/contractus.png)


# Contractus

Jamming around orders with API endpoints 🎸

## Environment Variables

To start the service locally, you need to export the current environment variables:

    export CONTRACTUS_POSTGRES_URL=<>

But, to play around the code and tests, it's possible with the command line ⤵️⤵️⤵️

## Command line
All commands are synthesized in the Makefile, to start the development environment, just run:

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

## Logs in production
    `heroku logs --tail -a contractus`

## API documentation
[API Docs](api/docs/)

## Attention points 
 - For a while, the integration-tests just ran locally not in CI, this increased the time to ship code 🚀
 - We already don't have a way to paginate transactions. 😔
 - To publish images and new releases, for now, the only way is using the command line, isn't automate by CI yet 😔
