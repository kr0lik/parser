# WIP: parser

chromedriver version 117

## Usage

Add `.env` from `.env.dist`

### Ozon

Add `app/config/ozon/filter.yaml` from `app/config/ozon/filter.yaml.dist`

Run `collectCategories`
```
$ docker-compose run --rm app ./bin/ozon collectCategories
```
and then `collectProducts`
```
$ docker-compose run --rm app ./bin/ozon collectProducts
```

### Wildberries

Add `app/config/wildberries/filter.yaml` from `app/config/wildberries/filter.yaml.dist`

Run `collectCategories`
```
$ docker-compose run --rm app ./bin/wildberries collectCategories
```
and then `collectProducts`
```
$ docker-compose run --rm app ./bin/wildberries collectProducts
```

## Dev
### Local Chrome

Add to build `-X 'parser/internal/infrastructure/selenium.devMode=y'` for use chrome with graphic interface:
```
$ go build -X 'parser/internal/infrastructure/selenium.devMode=y'"
```

### Docker Chrome
Start selenium in docker
```
$ docker-compose up selenium
```
and add in .env 
```
SELENIUM_HOST=selenium
```
and run

###### Debug
```
$ docker compose build app
$ docker run -p 2345:2345 --network parser_default --env-file=.env parser-app dlv --listen=:2345 --headless=true --log=true --api-version=2 --accept-multiclient exec ./bin/ozon <cmd>
or
$ docker run -p 2345:2345 --network parser_default --env-file=.env parser-app dlv --listen=:2345 --headless=true --log=true --api-version=2 --accept-multiclient exec ./bin/wildberries <cmd>
```
