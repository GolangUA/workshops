# Shop API

Let's create API of some common online shop:

## DB entity

- item [id, alias, title, desc, price, pictures, count]
- category [id, parent_id, alias, title, desc, logo]
- orders [id, sum, status, created, updated]

Each item belongs to the one category

## API endpoints

Client zone(/):
GET /items - return list of the items
GET /items/{itemID} - return one item
GET /categories - list of categories
GET /categories/{categoryID} - get category info
POST /orders - create an order
PUT /orders/{orderID}/item - add item to the order (by id and count)

GET /assets - for static files like pictures, js, css

## CLI

- import category {filePath}
- import logos {filePath}
- import items {filePath}
- import pictures {filePath}
- update order {orderID} {status}
- export orders {from-date} {to-date} {status}(optional)

## tools

[awsome-go](https://awesome-go.com/)

[docker](https://docs.docker.com/get-started/)
[docker-compose](https://docs.docker.com/compose/)

[Buffalo](https://gobuffalo.io/en)

[A simple, fast, and fun package for building command line apps in Go](https://github.com/urfave/cli)

[go-swagger](github.com/go-swagger/go-swagger)
[example of todo list server](https://goswagger.io/tutorial/todo-list)

[table tests generator](https://github.com/cweill/gotests)
[Mock generator](https://github.com/golang/mock)
[generator of fake data](https://github.com/icrowley/fake)

## Additional

### Authorization

### Admine zone API

Admin zone(/admin):
POST /items - create new item
PUT /items/{itemID} - update item info
DELETE /items/{itemID} - delete item
POST /items/{itemID}/picture - add picture to the item

POST /categories - add new category
PUT /categories/{categoryID} - update category info
DELETE /categories/{categoryID} - delete category info

GET /orders - list of orders
PUT /orders/{orderID} - update order info

POST /import/categories - upload list of categories (json)
POST /import/items - upload list of items (json)
POST /import/pictures - upload archive of pictures for existing items (zip, alias suffix matching)
GET /exoprt/orders - download list of orders