
# basket service
Basket service is a microservice that allows to generate a basket, add items to the basket, remove items from the basket
and checkout the basket.

This is an example of how you may give instructions on setting up your project locally. To get a local copy up and running follow these simple example steps.

## Installation

- Clone repo
```
git clone https://github.com/erdemcemal/basket-service.git
```
- Install dependencies
```
go mod download
```

## Usage

Project contains a docker-compose file. This file will create two container. One is postgresql and other one is basket-service api with alpine versions.

```
docker-compose up
```

If you have installed Taskfile on your machine you can also use taskfile comamnds to run project.
```
task run // build project and run the project
task test // execute unit test in the project
task build // build the project
```

````
localhost:8080 // basket service api running
````

> **_NOTE:_**  There is no need to add any products in the database. This is done automatically when you run the project. 
> Every time when you run the project migrations are executed. If there is no products in the database, they are added.

There are 7 endpoints available in the project. 

For "/alive" and "/products" endpoints there is no need to authenticate. For other endpoints you need to send a valid user_id in the header. For example in the header;
    
    ````
    curl --location --request GET 'http://localhost:8080/api/v1/basket' \
    --header 'user_id: 7f6c43bc-14a2-4b3a-898c-ae27a1d41b8d'
    ````
- /alive //returns 200 if service is running
```
http://localhost:8080/alive // check if service is running
```
- /api/v1/products //returns list of products
```
http://localhost:8080/api/v1/products // get list of products
```
- /api/v1/basket //returns basket related to user_id. If basket is not found, it will generate a new basket for the user_id and return it.
```
  curl --location --request GET 'http://localhost:8080/api/v1/basket' \
  --header 'user_id: 7f6c43bc-14a2-4b3a-898c-ae27a1d41b8d'
```
- /api/v1/basket // add item to the basket. If item is not found or item quantity is less than one, it will throw an error.
```
  curl --location --request POST 'http://localhost:8080/api/v1/basket' \
  --header 'Content-Type: application/json' \
  --data-raw '{
        "product_id": "9f1c3bb5-909f-4ccc-a77f-913bb398abc4",
        "quantity": 5
    }
    '
```

- /api/v1/basket/{productId} // remove item from the basket. If item is not found it will throw an error.
```
 curl --location --request DELETE 'http://localhost:8080/api/v1/basket/9f1c3bb5-909f-4ccc-a77f-913bb398abc4' \
 --header 'user_id: 7f6c43bc-14a2-4b3a-898c-ae27a1d41b8d'
```

- /api/v1/basket // update item quantity in the basket. If item is not found or item quantity is less than one, it will throw an error.
```
  curl --location --request PUT 'http://localhost:8080/api/v1/basket' \
    --header 'user_id: 7f6c43bc-14a2-4b3a-898c-ae27a1d41b8d' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "product_id": "9f1c3bb5-909f-4ccc-a77f-913bb398abc4",
        "quantity": 4
    }'
```

- /api/v1/basket/checkout // checkout the basket. No need any payment information.
```
  curl --location --request GET 'http://localhost:8080/api/v1/basket/checkout' \
    --header 'user_id: 7f6c43bc-14a2-4b3a-898c-ae27a1d41b8d'
```

## Campaign Engine (Discount apply on basket)

There are 3 different rules available for campaign engine. Only highest campaign will be applied on the basket.



### 1. Same product rule
If there are more than 3 items of the same product, then fourth and subsequent ones would have %8 off discount.

---
### 2. Every fourth order rule
Every fourth order whose total is more than given amount may have discount
depending on products. Products whose VAT is %1 don’t have any discount but products whose VAT is %8 and %18 have discount of %10 and %15 respectively.

---
### 3. Purchase amount rule
If the customer made purchase which is more than given amount in a month then all subsequent purchases should have %10 off.

---

To apply discount rules on basket, you need to specify the "given amount" in the docker-compose file. 

In the docker-compose file, you can change the "given amount" by specifying the "GIVEN_AMOUNT" environment variables in the docker-compose file under the "services/api" section.


---

## Tech Stack
- GoLang
- Postgresql
- Docker
- Docker Compose
- Taskfile

Applied practices:
- gorm: ORM for GoLang
- validator: validation library for GoLang. It is used to validate input parameters.
- middleware: add logging, authentication and json middleware to the project.
- gorilla-mux: used to handle http requests.
- handler -> service -> repository implementation: used to implement business logic.
- custom error handling: used to handle errors.
- dto: data transfer object used to transfer data between services.
- logrus: used to log errors.