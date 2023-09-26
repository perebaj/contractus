# contractus

![Logo](/assets/blackbackground.png#gh-dark-mode-only)
![Logo](/assets/whitebackgroun.png#gh-light-mode-only)

![image](assets/contractus.png)


# Contractus

Jamming around orders with API endpoints 🎸

## Environment Variables

To start the service you need to export the current environment variables:

    export CONTRACTUS_POSTGRES_URL=<>

## Command line
All commands should be runned within a container environment:
    
    make dev/start
    make dev

Then you can run:
    
    The integration tests: make integration-test
    The Unit tests: make test
    The lint: make lint

## Ship a new version
    `make image/publish`
    `heroky container:release web -a contractus`

## Logs in production

    `heroku logs --tail -a contractus`

## API documentation



