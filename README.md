# sagara-test-BE

the goal of these task are to create CRUD API, auth JWT and upload picture 

## setup
1. use dockerfile to create docker image
2. type ```make docker.fiber``` to build and run
3. Hit the Server to test Health `localhost:5000/api/v1/health`
4. let's Rock !! 🚀

### Tasks 
We define routes for handling operations:
prefix ```/api/v1```

| Method        | Route                  | Action                                              |
|---------------|------------------------|-----------------------------------------------------|
| GET           | /login                 | create token JWT                                    |
| POST          | /book                  | create book                                         |
| POST          | /book/:id/upload-image | upload image and update data                        |
| GET           | /books                 | get all books                                       |
| GET           | /book/:id              | get book by id                                      |
| PUT           | /book                  | update book data                                    |
| DELETE        | /book                  | delete book data                                    |

Access API via ```http://localhost:5000/api/v1{route}```


1. GET ```/login```

Response:
status code: 200
```
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjkyMzU5MjF9.C0mI1Oc_U4E5A2l_qXhToFaw3epgHn1Jj2S1K2EnTdQ",
    "error": false,
    "msg": "success create token"
}
```

2. POST ```/book ```

Authorization: Bearer {token} 

Request Body: 
```
{
    "title": "da vinci code",
    "author": "dan brown",
    "book_attrs": {
        "description": "bercerita tentang artefak bunda maria",
        "rating": 9
    }
}
```

Response:
status code : 201
```
{
    "book": {
        "id": "fe7d6284-6c62-4b8f-a745-32be28eadb66",
        "created_at": "2021-08-18T04:19:32.715601+07:00",
        "updated_at": "0001-01-01T00:00:00Z",
        "title": "da vinci code",
        "author": "dan brown",
        "book_status": 1,
        "book_attrs": {
            "picture": "",
            "description": "bercerita tentang artefak bunda maria",
            "rating": 9
        }
    },
    "error": false,
    "msg": "success create data"
}
```

3. POST ```/book/:id/upload-image ```

Authorization: Bearer {token} 

Headers : 
    content-type: multipart/form-data

form-data: 
```
key : image <file>
value: photo.jpg
```

Response:
status code: 201
```
{
    "error": false,
    "msg": "success upload picture"
}
```

4. GET ```/books```

Authorization: Bearer {token} 

Response:
status code: 200
```
{
    "books": [
        {
            "id": "10669dcf-75af-481a-870e-e02771a8d0ca",
            "created_at": "2021-08-18T03:05:47.011951+07:00",
            "updated_at": "2021-08-18T03:52:53.579299Z",
            "title": "laskar pelangi",
            "author": "batman",
            "book_status": 1,
            "book_attrs": {
                "picture": "https://res.cloudinary.com/dfftzrpvh/image/upload/v1629228794/my_image.jpg",
                "description": "bercerita seorang anak di pulai belitong",
                "rating": 8
            }
        },
        {
            "id": "fe7d6284-6c62-4b8f-a745-32be28eadb66",
            "created_at": "2021-08-18T04:19:32.715601+07:00",
            "updated_at": "2021-08-18T04:21:25.275764Z",
            "title": "da vinci code",
            "author": "dan brown",
            "book_status": 1,
            "book_attrs": {
                "picture": "https://res.cloudinary.com/dfftzrpvh/image/upload/v1629228794/my_image.jpg",
                "description": "bercerita tentang artefak bunda maria",
                "rating": 9
            }
        }
    ],
    "count": 2,
    "error": false,
    "msg": "success get all book"
}
```

5. GET ```/books/:id```

Authorization: Bearer {token} 

Response:
status code: 200
```
{
    "book": {
        "id": "fe7d6284-6c62-4b8f-a745-32be28eadb66",
        "created_at": "2021-08-18T04:19:32.715601+07:00",
        "updated_at": "2021-08-18T04:21:25.275764Z",
        "title": "da vinci code",
        "author": "dan brown",
        "book_status": 1,
        "book_attrs": {
            "picture": "https://res.cloudinary.com/dfftzrpvh/image/upload/v1629228794/my_image.jpg",
            "description": "bercerita tentang artefak bunda maria",
            "rating": 9
        }
    },
    "error": false,
    "msg": "success get book by id"
}
```

6. PUT ```/book```

Authorization: Bearer {token} 

Request : 
```
{
    "id":"fe7d6284-6c62-4b8f-a745-32be28eadb66",
    "author": "mario" 
}
```

Response:
status code: 201
```
{
    "error": false,
    "msg": "success update data"
}
```

7. DELETE ```/book```
Authorization: Bearer {token} 

Request : 
```
{
    "id":"ac3886b7-396b-4f5a-91b5-30767f5bc07d"
}
```

Response:
status code: 200
```
{
    "error": false,
    "msg": "success delete data with id : ac3886b7-396b-4f5a-91b5-30767f5bc07d"
}
```

### Tech Stack
* [Golang] - programming language
* [Fiber] - web framework with zero memory allocation and performance
* [ElephantSQL] -  browser tool for SQL queries where you can create, read, update and delete data directly from your web browser
* [Cloudinary] - like AWS S3 storage
* [JsonWebToken] - Authorization and Authentication 


[Golang]: <https://golang.org/>
[Fiber]: <https://github.com/gofiber/fiber/>
[ElephantSQL]: <https://www.elephantsql.com/>
[Cloudinary]:<https://cloudinary.com//>
[JsonWebToken]: <https://jwt.io/>

*Notes
why i choose elephantSQL ? for table migration need to install ``` golang-migrate```,  better i use database online ``` ElephantSQL``` so no need to setup the apps but . . . . because its database online resulting in slower performance due to latency and network, please understand.
