# Go-Chi-JWT-REST-API-Demo

This is a simple Go app to show demonstration of JWT Authentication in REST API.

In the code I use [chi](https://github.com/go-chi/chi) and [jwtauth](https://github.com/go-chi/jwtauth) for the demonstration.

The `/articles` API endpoints are the protected endpoints which requires authentication with JWT token to be passed in the request header.

To acquire the JWT token, calls the `/login` API endpoint with the supplied user id ("admin") and password ("123456"). The JWT token will be returned in the response. e.g. `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY1NDQyMTgsImlkIjoiYWRtaW4ifQ.vt_7g7Cg_epngB1ZCnuuWl41bExunNOJXI8fT5mKO7U"}`

After that, you would pass the token in the `Authorization` header, e.g. `'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY1NDQyMTgsImlkIjoiYWRtaW4ifQ.vt_7g7Cg_epngB1ZCnuuWl41bExunNOJXI8fT5mKO7U'`

For more details, check the **How to use the app** section below. Also check the `main.go` source code file to see how the authentication is implemented.

## Requirements

 * Go

## How to install dependencies

	$ cd Go-Chi-JWT-REST-API-Demo
	$ go get ./...

## How to run the app

	$ go run main.go

The app should run by listening on localhost:3000.

## How to use the app

1) By using curl command, try to call the `GET /articles` API first without JWT token:

    ```
    $ curl -XGET 'http://localhost:3000/articles'
    ```

    It should return response:

    `{"message":"Not authorized"}`

2) Next, call the `POST /login` API to acquire the token:

    ```
    $ curl -XPOST 'http://localhost:3000/login' \
    --header 'Content-Type: application/json' \
    --data '{
        "userId": "admin",
        "password": "123456"
    }'
    ```

    It should return the token in the response, e.g.

    `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY1NDQyMTgsImlkIjoiYWRtaW4ifQ.vt_7g7Cg_epngB1ZCnuuWl41bExunNOJXI8fT5mKO7U"}`

    Note that the token is set to expire in 10 minutes. Check the code in `main.go` to see how the expiry is set.

3) Pass the token to the `GET /articles` API call:

    ```
    $ curl -XGET 'http://localhost:3000/articles' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY1NDQyMTgsImlkIjoiYWRtaW4ifQ.vt_7g7Cg_epngB1ZCnuuWl41bExunNOJXI8fT5mKO7U'
    ```

    It should return response:

    `[{"id":"1","title":"Hello 1","desc":"Article Description 1","content":"Article Content 1"},{"id":"2","title":"Hello 2","desc":"Article Description 2","content":"Article Content 2"},{"id":"3","title":"Hello 3","desc":"Article Description 3","content":"Article Content 3"},{"id":"4","title":"Hello 4","desc":"Article Description 4","content":"Article Content 4"},{"id":"5","title":"Hello 5","desc":"Article Description 5","content":"Article Content 5"}]`

    If you call the API after the token had expired, you should receive this response:

    `{"message":"Not authorized"}`
