# Blog API

## Prompt

Build a backend API that allows users to create, edit, publish, and retrieve blog posts. This is an
API-only project — no frontend required. Please timebox this assignment and do not spend
more than 4 hours on it.

### Features & Requirements

#### Core Functionality

- RESTful API
- User management can be mocked or hardcoded
- CRUD for blog posts
- Only allow the author to update/delete their own posts
- Published posts should be readable by any user
- Drafts should only be accessible to their author
- Data storage - database is preferred

#### Technical Expectations

- Written in Go or NodeJS
- Well-structured code with clear separation of concerns
- Basic test coverage
- Linting, formatting, and adherence to idiomatic practices

##### Bonus

- Pagination support
- Admin-only endpoint to delete any post

#### Deliverable

- GitHub repository link (or zip)
- README file with:
  - How to run the project
  - Explanation of the design decisions
  - Any tradeoffs or limitations
- Optional: Example curl commands or Postman file

## Running this project

TODO

## Design Decisions

### Identity and Access

Implicitly, there is a need for identity for users that are creating posts that will be associated to them. Our API should acertain that the authors are who they say they are when they are creating posts on our site. Further, there are a number of constraints on which users are authorized to perform certain actions:

- Only allow the author to update/delete their own posts
- Published posts should be readable by any user
- Drafts should only be accessible to their author
- Admin-only endpoint to delete any post

This points me to a need to establish a basic level of authentication and authorization. JWT tokens help handle both of these by allowing my service to sign a set of claims about my user that they can in turn use from the client side to make requests. Doing this requires:

- A login endpoint to provide the user with tokens
- A JWT implementation on the backend
- Auth middleware on my endpoints to help unwrap user identity on each request

My user data model contains one field to indicate authorization (`is_admin`) but as the user model and permissions expand this approach will be difficult to scale. Each modification to what a user can do would require an update the user model and underlying DB. A more robust solution would be to set up a one to many relationship between a user and their roles or permissions to establish more extensible access control.

### Storage Concerns

Most of the information this API needs to store and retieve are basic metadata fields. The article title, author name, and description can all be represented by relatively short text fields. For the sake of time, I am going to use PostgreSQL to set up users and relationships due to my familiarity with the technology. Arguably, depending on how this information could be modelled as unstructured data in a NoSQL DB technology, especially is the articles themselves required a JSON structure to model rich text or formatting declarations.

## Tradeoffs and Limitations

### Rich Text Editing

My service is supporting only a simple title, byline and text blob as a blog post. Ideally, a user would be able store richer data that allows for formatting of the text (size, color, emphasis) and subheadings. There are a few options for storing the article as more structured data.

- **HTML** - This being the lingua franca of the browser makes it a an easy choice for rendering purposes. This needs to be handled carefully though as it brings with it injection vulnerabilities if not properly sanitized. It also is opinionated about what kind of clients should be presenting the articles - it's possible we want to support an iOS or Android app down the road which may make the choice to use HTML short sighted.
- **Markdown** - This may be a nice option if your author user base is already familiar with the syntax. The resulting text blob would be easily searchable and a module for rendering MD into HTML or a different presentation format would separate the presentation details from the raw text content. The limitation is that the feature set for editing would be limited to what Markdown syntax supports.
- **JSON** - A JSON structure could be used to set up a custom set of directives / attributes for styling and is agnostic of any specific presentation format. The overhead here is setting up a client that knows how to convert what is in the editor into this proprietary structure and maintaining that API over time.

### Editing History

The data model I am proposing only supports a single snapshot of each blog. This will function well for a simple prototype, but in a full-fledged blog editing app I may want to support reviewing revisions. Similar to git, our service could support a log of revisions to capture the full history of the the blog post and offer authors more powerful tools for editing and reviewing the articles they are working on. This however, would require a more complex data model and more storage for each blog post.

### Search Functionality

Another ideal feature as the dataset of blog posts grows is a means for performing text or tag search against the blog posts. The current API only supports paginated search but from a UX standpoint a user would likely want to be able to filter for posts they are interested in.

## API Spec

This is my "top-down" method for modeling the problem.

### Login

#### Request

```http
POST /api/v1/login HTTP/1.1
Host: localhost:8080
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
```

```sh
curl --location --request POST 'http://localhost:8080/api/v1/login' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ='
```

#### Responses

##### 200 - User Login Accepted

```json
{
  "token": "{{jwt_token}}"
}
```

##### 401 - User Authentication Details Incorrect

```json
{
  "error": "user does not exist or wrong password provided"
}
```

### Posts

#### Create Post

##### Request

```http
POST /api/v1/posts HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer {jwt_token}
Content-Length: 150

{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "content": "Some really long string"
}
```

```sh
curl --location 'http://localhost:8080/api/v1/posts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {jwt_token}' \
--data '{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "content": "Some really long string"
}'
```

##### Responses

###### 201 - Created

```json
{
  "post_id": "57e88e7f-2974-45ef-8e6d-87ac81ad81c2"
}
```

###### 400 - Bad Request

Returned if

- required field is missing
- data constraint broken (ie: title duplicated)

```json
{
  "error": "required field 'title' is missing"
}
```

###### 401 - Unauthorized

Returned if:

- User does not exist or JWT is invalid

```json
{
  "error": "invalid credentials"
}
```

#### Fetch All Posts

Auth is not needed for this endpoint as there shouldn't be a restriction on read. This endpoint should not include drafts in the response.

##### Request

```http
GET /api/v1/posts HTTP/1.1
Host: localhost:8080
```

```curl
curl --location 'http://localhost:8080/api/v1/posts'
```

##### Responses

###### 200 - OK

```json
{
  "posts": [
    {
      "author": {
        "id": "f98e0378-b419-464d-875d-e75ec0124e4c",
        "name": "Kazuo Ishiguro"
      },
      "id": "6fb0e026-333c-49ff-965c-1615b30dad57",
      "title": "Klara and the Sun",
      "summary": "Some summary under N chars",
      "created_ts": "2025-06-24T21:53:44Z",
      "published_ts": "2025-06-24T21:53:44Z",
      "updated_ts": "2025-06-24T21:53:44Z"
    }
  ]
}
```

#### Fetch Post by ID

##### Request

```http
GET /api/v1/posts/:id HTTP/1.1
Host: localhost:8080
Authorization: Bearer {{jwt_token}}
```

```sh
curl --location 'http://localhost:8080/api/v1/posts/:id' \
--header 'Authorization: Bearer {{jwt_token}}'
```

##### Responses

###### 200 - OK

```json
{
  "post": {
    "author": {
      "id": "f98e0378-b419-464d-875d-e75ec0124e4c",
      "name": "Kazuo Ishiguro"
    },
    "id": "6fb0e026-333c-49ff-965c-1615b30dad57",
    "title": "Klara and the Sun",
    "summary": "Some summary under N chars",
    "content": "Some really long string",
    "status": "PUBLISHED",
    "created_ts": "2025-06-24T21:53:44Z",
    "published_ts": "2025-06-24T21:53:44Z",
    "updated_ts": "2025-06-24T21:53:44Z"
  }
}
```

###### 401 - Unauthorized

Returned if:

- User does not exist or JWT is otherwise invalid

```json
{
  "error": "invalid credentials"
}
```

###### 403 - Forbidden

Returned if:

- User does not have access to draft article

```json
{
  "error": "user does not have access to this blog post"
}
```

###### 404 - Not Found

Returned if:

- Article does not exist

```json
{
  "error": "blog post requested was not found"
}
```

#### Update Post

##### Request

```http
PUT /api/v1/posts/:id HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer {jwt_token}
Content-Length: 150

{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "content": "Some really long string"
}
```

```sh
curl --location --request PUT 'http://localhost:8080/api/v1/posts/:id' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {jwt_token}' \
--data '{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "content": "Some really long string"
}'
```

##### Responses

###### 200 - OK

```json
{}
```

###### 400 - Bad Request

Returned if

- required field is missing
- data constraint broken (ie: title duplicated)

```json
{
  "error": "required field 'title' is missing"
}
```

###### 401 - Unauthorized

Returned if:

- User does not exist or JWT is invalid

```json
{
  "error": "invalid credentials"
}
```

###### 403 - Forbidden

Returned if:

- User is not the author of the post

```json
{
  "error": "user cannot update this entity"
}
```

#### Delete Post

##### Request

```http
DELETE /api/v1/posts/:id HTTP/1.1
Host: localhost:8080
Authorization: Bearer {jwt_token}
```

```curl
curl --location --request DELETE 'http://localhost:8080/api/v1/posts/:id' \
--header 'Authorization: Bearer {jwt_token}'
```

##### Responses

###### 204 - No Content

###### 401 - Unauthorized

Returned if:

- User does not exist or JWT is invalid

```json
{
  "error": "invalid credentials"
}
```

###### 403 - Forbidden

Returned if:

- User is not the author of the post

```json
{
  "error": "user cannot update this entity"
}
```

#### Delete Post (Admin)

##### Request

```http
DELETE /api/v1/admin/posts/:id HTTP/1.1
Host: localhost:8080
Authorization: Bearer {jwt_token}
```

```curl
curl --location --request DELETE 'http://localhost:8080/api/v1/admin/posts/:id' \
--header 'Authorization: Bearer {jwt_token}'
```

##### Responses

###### 204 - No Content

###### 401 - Unauthorized

Returned if:

- User does not exist or JWT is invalid

```json
{
  "error": "invalid credentials"
}
```

###### 403 - Forbidden

Returned if:

- User is not an admin

```json
{
  "error": "user does not have admin access"
}
```

## Data Model

This is my "bottom-up" way of modelling the problem.

### Entities

#### Users

##### Attributes

| Field      | Data Type | Description                                         |
| ---------- | --------- | --------------------------------------------------- |
| `id`       | uuid      | unique identifier for the user                      |
| `username` | string    | login identity                                      |
| `name`     | string    |                                                     |
| `is_admin` | boolean   | indicates whether user has special admin privileges |

#### Posts

##### Attributes

| Field          | Data Type               |
| -------------- | ----------------------- |
| `id`           | uuid                    |
| `status`       | enum (PUBLISHED, DRAFT) |
| `title`        | string                  |
| `summary`      | string                  |
| `contents`     | string                  |
| `author_id`    | uuid                    |
| `created_ts`   | timestamp               |
| `published_ts` | timestamp               |
| `updated_ts`   | timestamp               |

## Miscellaneous Details

### JWT Token Structure

This JWT token will be central to my authentication and authorization strategy for the API. I will make a claim about the user's identity and admin status here.

```json
{
  "sub": "1234567890",
  "user_id": "1ecaf3dc-db60-468e-a404-04b7a7d521c1", // establish user identity
  "is_admin": true, // makes a claim about user authorization level
  "iat": 1516239022
}
```

### Site Users

The following is the list of site users mocked for usage.

| Username  | Password    | Is Admin |
| --------- | ----------- | -------- |
| kishiguro | hailsham    | false    |
| dsedaris  | emeraldIsle | false    |
| admin     | password    | true     |
