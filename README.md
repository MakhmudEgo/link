# Link
####Укоротитель ссылок, предоставляет API к созданию коротких ссылок.
#Запуск:
####Cервис работает на postgres или redis
##postgres:
		make runPostgres
##redis:
		make runRedis


# API:
##### 	http://localhost:8080/link
## Примеры запросов:
### POST http://localhost:8080/link
### body
		{
			"url": "https://ozon.dev/internship"
		}
## Примеры ответов:
		{
			"url": "http://localhost:8080/XxXXxXXxXX",
			"error": false,
			"description": "Created"
		}

		{
			"url": "http://localhost:8080/XxXXxXXxXX",
			"error": false,
			"description": "Success"
		}

		{
			"url": "",
			"error": true,
			"description": "Server URL Are Not Supported"
		}

### GET http://localhost:8080/link
### body
		{
			"url": "http://localhost:8080/xxxxxxxxxX"
		}
## Примеры ответов:
		{
			"url": "https://ozon.dev/internship",
			"error": false,
			"description": "Success"
		}

		{
			"url": "",
			"error": true,
			"description": "No Server URL Are Not Supported"
		}
