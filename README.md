# Cozy Take Home Test

Thanks for you interest in working with us here at Cozy! This 
is a small take home assessment to test your proficieancy with 
building backend web services.

## Setup

You will need [Go](https://go.dev/) and [Docker](https://www.docker.com/) to run this project.

Here are the steps to get the project setup:
- Run `docker compose up`
  - You can specify the `-d` flag if you want to start it in the background.
- Run `go run main.go seed`
  - This will setup the DB schema and seed it with dummy data.
- Run `go run main.go`
  - This will start up a basic RESTful server with a single endpoint.

We will be using Go with the Chi router and 
Postgres (Running in the Docker container) for this. The schema 
for the SQL DB can be found in [schema.sql](./sql/schema.sql) and 
a basic Chi setup can be found in [api/routes.go](./api/routes.go). 

## assessment

I want you to create an RESTful API with the following endpoints:

- `/posts`
  - This endpoint should return a paginated list of posts.
  - The return should be a list of objects with
    - post data
    - author data
    - the number of likes
  - It should take an optional `user` query param and if passed should return
    if that user liked the post.
- `/posts/:id`
  - This endpoint should return a specific post by its ID.
  - The return should be an object with
    - post data
    - author data
    - the number of likes
  - It should take an optional `user` query param and if passed should return
    if that user liked the post.
- `/posts/:id/likes`
  - This endpoint should return a paginated list of the users who liked a post.
  - The return should be a list of objects with
    - like date
    - user data
- `/users/:id`
  - This endpoint should return a specific user by their ID.
  - The return should be an object with
    - user data
    - 5 latest posts by this user

## Questions

If you have any questions or issues, please reach out to me at [graham@gvasquez.dev](mailto:graham@gvasquez.dev).
