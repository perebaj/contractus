# contractus

![Logo](/assets/blackbackground.png#gh-dark-mode-only)
![Logo](/assets/whitebackgroun.png#gh-light-mode-only)

![image](assets/contractus.png)


# Contractus

Jamming around orders with API endpoints 🎸

## Environment Variables

If you want to run this project locally you must set up these environment variables:

    CONTRACTUS_GOOGLE_CLIENT_ID
    CONTRACTUS_GOOGLE_CLIENT_SECRET

This can be done in your terminal, where you export the variables, and then in the same terminal, run the service startup command `make dev/start`

**Note:** To create your own Google Client for OAuth2.0, access the following [link](https://console.cloud.google.com/apis/credentials), without these secrets, the service won't run locally, but it's possible to access it using the link present in the repository details on GitHub

## Get Started

To start the service locally, you can type `make dev/start` and after that you can use the docker container IP to play around the routes, `make ip`

Request example:

    curl (make ip)

Or just accessing `http:/localhost:8080`

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

## It's also good to know 
- You can log in to the service using your Google account. No bureaucracy to 🎸
- Structured logs all the way 🥸
- Deployed to the open sea through Heroku 🌊 (Check the repository details to access the link)
---
- For a while, the integration-tests just ran locally  not in CI 😔, this increased the time to 🚀 code;
- We don't have a way to paginate transactions yet; 😔
- To publish images and new releases, for now, the only way is using the command line, isn't automate by CI yet; 😔
- The infra isn't automated by the power of the IAC yet.😔 Button engineer only. 🔘✅

