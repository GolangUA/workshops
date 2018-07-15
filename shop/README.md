# Shop API

Let's create API for some common online shop:

## DB entity

* item [id, alias, title, desc, price, pictures, count]
* category [id, parent_id, alias, title, desc, logo]
* orders [id, sum, status, created, updated]

Each item belongs to the one category

## API endpoints

Client zone(/):

* GET /items - return list of the items
* GET /items/{itemID} - return one item
* GET /categories - list of categories
* GET /categories/{categoryID} - get category info
* POST /orders - create an order
* PUT /orders/{orderID}/item - add item to the order (by id and count)

* GET /assets - for static files like pictures, js, css

## tools

* [awsome-go](https://awesome-go.com/)

* [docker](https://docs.docker.com/get-started/)
* [docker-compose](https://docs.docker.com/compose/)

* [Buffalo](https://gobuffalo.io/en)

* [A simple, fast, and fun package for building command line apps in Go](https://github.com/urfave/cli)

* [go-swagger](github.com/go-swagger/go-swagger)
* [example of todo list server](https://goswagger.io/tutorial/todo-list)

* [table tests generator](https://github.com/cweill/gotests)
* [Mock generator](https://github.com/golang/mock)
* [generator of fake data](https://github.com/icrowley/fake)

## Steps

### 0. Install Buffalo

```
go get -u -v github.com/gobuffalo/buffalo/buffalo
```

### 1. Generate API

```
buffalo new api --api --vcs none --docker none
```

for launch API localy (http://localhost:3000):
```
buffalo dev
```

four build docker image use this:
```
cd ./api
GOOS=linux buffalo build -k -o ./app
cd ../
docker build -t shop-api .
```

Note: use last 2 commands for rebuild application

### 2. Generate Routes

```
buffalo generate actions items list index --method GET --skip-template
buffalo generate actions categories list index --method GET --skip-template
buffalo generate actions orders create --method POST --skip-template
buffalo generate actions orders update --method PUT --skip-template
```

Note: You should change [render](https://godoc.org/github.com/gobuffalo/buffalo/render#Engine.JSON) from HTML to JSON manually at each handler

Note: Remember about API prefix `/v1` and other routes conventions from Swagger

Validation: `buffalo routes` or inside container `api`: `./app task routes`

### 4. Generate Models

```
soda generate model item alias:string title:string desc:string pictures:string price:int count:int category_id:uuid
soda generate model category alias:string title:string desc:string logo:string parent_id:uuid
soda generate model order id:int status:string sum:int
soda generate model ordered id:int order_id:int item_id:uuid item_cnt:int item_sum:int
```

Note: You may change structure if have your own vision at data structure.

Note: for integrate with DB use docker-compose:
```
docker-compose up -d # you should have shop-api image
```

on any changes at codebase you should recreate docker image each time and after that:
```
docker-compose up -d api
```

for logs observation:
```
docker-compose logs -f api # without flag -f - just output all existed logs and stop
```

for migration execution:
```
docker-compose exec api /bin/sh
./app migrate
exit
```



### 5. Create Seeds task

Based on [fake](https://github.com/icrowley/fake) generate categories (two levels: 5-15 first level with 5-20 on each), items for 50-150 items on each category

Note: update `database.yml` with proper DB access creds: `host=db user=postgres password=postgres database=shop`

Note: for execution go inside container `docker-compose exec api /bin/sh` and run task: `./app task seed`

### 6. Implement handler logic

Write business logic for selecting, filtering of categories, items, create and update orders

### 7. Implement tooling for shop admin

* import category {filePath}, simple csv format
* import logos {filePath}, should unpach zip file and update categories with corrisponding pictures (match by alias as suffics), store referenece to the files, handler for static files should be added
* import items {filePath}, simple csv format
* import pictures {filePath}, should unpach zip file and update items with corrisponding pictures (match by alias as suffics), store JSON string  with list of referenece to the files, handler for static files should be added
* update order {orderID} {status}, use data from csv file
* export orders {from-date} {to-date} {status}(optional), generate csv file with orders and suppor for filter by date and status

## Additional

### Authorization

### Admine zone API

Admin zone(/admin):
* POST /items - create new item
* PUT /items/{itemID} - update item info
* DELETE /items/{itemID} - delete item
* POST /items/{itemID}/picture - add picture to the item

* POST /categories - add new category
* PUT /categories/{categoryID} - update category info
* DELETE /categories/{categoryID} - delete category info

* GET /orders - list of orders
* PUT /orders/{orderID} - update order info

* POST /import/categories - upload list of categories (json)
* POST /import/items - upload list of items (json)
* POST /import/pictures - upload archive of pictures for existing items (zip, alias suffix matching)
* GET /export/orders - download list of orders
