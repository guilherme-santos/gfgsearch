# Global Fashion Group Search

Search engine written in Go using ElasticSearch as database.

## Running

You have two different ways to run this app. You can use docker-compose or run manually from your machine (go is required).

#### Docker-Compose 

First of all you need to create **gfgsearch** docker image, you can do that using following command:

```
docker-compose build
```

After that you can run **elasticsearch** and **gfgsearch** in background using:

```
docker-compose up -d
```

By default elasticsearch port is not exposed in the localhost, only the app is exposed. If you need to expose it
copy the file `docker-compose.override.yml.dist` to `docker-compose.override.yml` and run docker-compose
again. After that elasticsearch will be available in the port 9200.

**If you need to check or change the port that the app is running you can do it in the field `services.gfgsearch.ports` of file `docker-compose.yml`.**

## Manually

Before run the app you need to be sure that you have go installed in your machine you can do a quick check typing `go version`. You also need to have `dep` (go dependency management tool), if you don't have it you can download it [here](https://github.com/golang/dep/releases).

If you already have all needed tools, you need to install all project dependecies with the following command:

```
dep ensure
```

Before we run the app we need to run elasticsearch (if you don't have one running yet). By default elasticsearch is not exposing its port, so you can copy the file `docker-compose.override.yml.dist` to `docker-compose.override.yml` and type:

```
docker-compose up -d elasticsearch
```

Now you're ready to run the app, just type:

```
make run
```

**If you need to check or change the port that the app is running or elasticsearch address you can do it in the file `.env`. In case you don't have the file just copy `.env.dist` to `.env`.**

## Populating

Before you start to use the app you need to have some data available, if you're running docker-compose you need to copy a file with some data to inside of the container and then execute the *populate* command, you can do that as shown:

```
docker cp elasticsearch/testdata/products.json gfgsearch_gfgsearch_1:/tmp/products.json
docker-compose exec gfgsearch gfgsearch -populate /tmp/products.json
```

In case you're running manually, it's way simpler, just type:
```
make populate
```

## Accessing

Now that you have the app running and with some data you can access using (by default) http://localhost:8080/v1/search/products

It'll ask you for username and password, case you haven't setted it, check/change in the `docker-compose.yml` or `.env` file, depending how you have ran the app. By default username is **gfg** and password is **search**.

The API provides search, filtering, pagination and sorting. You can also combine all these features together.

| Features |  |
| - | - |
| Search | `/v1/search/products?q=shirt` |
| Filter | `/v1/search/products?filter=brand:adidas,price:900` |
| Page | `/v1/search/products?page=3` |
| Result per page | `/v1/search/products?per_page=40` |
| Sort asceding | `/v1/search/products?sort=price` |
| Sort desceding | `/v1/search/products?sort=-price` |

## Testing

You can easily run the tests typing `make test` but also are available some integration tests that can be run using `make integration-test`. Remember that to run integration tests you need elasticsearch to be running and the port should be exposed in the host machine as shown before.

Case you want to run a only a specific test you can do it, like:

```
make integration-test TESTCASE=TestSearch_SortBy
```

The previous command will run the test `TestSearch_SortByPriceAsc` and `TestSearch_SortByPriceDesc`.
