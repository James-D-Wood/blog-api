# Manual Testing Cases

Below is a sampling of how to manually test different functionalities. This is not exhaustive list but covers the requirements in the prompt.

- Only allow the author to update/delete their own posts
- Published posts should be readable by any user
- Drafts should only be accessible to their author

## User Log In

- Admin Username: `admin`
- Other Available Usernames: `dsedaris`, `kishiguro`

### Admin User

Request with Basic Auth for Admin User

```sh
curl --location --request POST 'http://localhost:8080/api/v1/login' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ='
```

Response: 200

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3YS04OGQxLWFiMjUzOGJmYzM3YSJ9.mlLNE0sv0YwiKrMYLfy5WmkW9KadVvsjBM-x4WNE7M4"
}
```

Decoded JWT

```json
{
  "is_admin": true,
  "user_id": "0197aad2-96f6-737a-88d1-ab2538bfc37a"
}
```

### Regular User

Request with Basic Auth for Standard User (`dsedaris`)

```sh
curl --location --request POST 'http://localhost:8080/api/v1/login' \
--header 'Authorization: Basic ZHNlZGFyaXM6cGFzc3dvcmQ='
```

Response: 200

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4"
}
```

Decoded JWT

```json
{
  "is_admin": false,
  "user_id": "0197aad2-96f6-7376-8945-18690000ed92"
}
```

### Unregistered User

```sh
curl --location --request POST 'http://localhost:8080/api/v1/login' \
--header 'Authorization: Basic bm90QVVzZXI6cGFzc3dvcmQ='
```

Response: 401

```json
{
  "error": "user does not exist or wrong password provided"
}
```

## Create Blog

### User Has Invalid Token

Here, I tampered with the claims to show how the signature is invalidated event if the JWT is correctly encoded.

```sh
curl --location 'http://localhost:8080/api/v1/posts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3Ni04OTQ1LTE4NjkwMDAwZWQ5MiJ9.PxwaOiijuSAvdzXenp9NIjLBM40K5m2Ie4NpCEGI2oY' \
--data '{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "contents": "Some really long string"
}'
```

Status: 401

```json
{
  "error": "could not authenticate user"
}
```

### User Has No Token

```sh
curl --location 'http://localhost:8080/api/v1/posts' \
--header 'Content-Type: application/json' \
--data '{
    "title": "My riveting blog post",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "contents": "Some really long string"
}'
```

Response: 401

```json
{
  "error": "could not authenticate user - Authorization header missing"
}
```

### User Has Valid Token

```sh
curl --location 'http://localhost:8080/api/v1/posts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4' \
--data '{
    "title": "My riveting blog post title",
    "status": "DRAFT",
    "summary": "Some summary under N chars",
    "contents": "Some really long string"
}'
```

Response: 201

```json
{
  "post": {
    "id": "0197aae1-8c00-7e25-80e6-81fc72c7d5f8",
    "status": "DRAFT",
    "title": "My riveting blog post title",
    "summary": "Some summary under N chars",
    "contents": "Some really long string",
    "author_id": "0197aad2-96f6-7376-8945-18690000ed92",
    "created_ts": "2025-06-25T23:16:37-07:00",
    "published_ts": "",
    "updated_ts": "2025-06-25T23:16:37-07:00"
  }
}
```

## Read Blog

### Owner Can Read Draft

```sh
curl --location 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4'
```

Response: 200

```json
{
  "post": {
    "id": "0197aae1-8c00-7e25-80e6-81fc72c7d5f8",
    "status": "DRAFT",
    "title": "My riveting blog post title",
    "summary": "Some summary under N chars",
    "contents": "Some really long string",
    "author_id": "0197aad2-96f6-7376-8945-18690000ed92",
    "created_ts": "2025-06-25T23:16:37-07:00",
    "published_ts": "",
    "updated_ts": "2025-06-25T23:16:37-07:00"
  }
}
```

### Other Users Cannot Read Draft

```sh
curl --location 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3YS04OGQxLWFiMjUzOGJmYzM3YSJ9.mlLNE0sv0YwiKrMYLfy5WmkW9KadVvsjBM-x4WNE7M4'
```

Response: 403

```json
{
  "error": "user not authorized to view this post"
}
```

## Update Blog

### Owner Can Update

```sh
curl --location --request PUT 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4' \
--data '{
    "status": "PUBLISHED",
    "title": "My riveting blog post with a new title",
    "summary": "Some summary under N chars",
    "contents": "Some really, really long string"
}'
```

Response: 200

```json
{
  "post": {
    "id": "0197aae1-8c00-7e25-80e6-81fc72c7d5f8",
    "status": "PUBLISHED",
    "title": "My riveting blog post with a new title",
    "summary": "Some summary under N chars",
    "contents": "Some really, really long string",
    "author_id": "0197aad2-96f6-7376-8945-18690000ed92",
    "created_ts": "2025-06-25T23:16:37-07:00",
    "published_ts": "",
    "updated_ts": "2025-06-25T23:24:29-07:00"
  }
}
```

### Other Users Cannot

```sh
curl --location --request PUT 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3YS04OGQxLWFiMjUzOGJmYzM3YSJ9.mlLNE0sv0YwiKrMYLfy5WmkW9KadVvsjBM-x4WNE7M4' \
--data '{
    "status": "PUBLISHED",
    "title": "My riveting blog post with a new title",
    "summary": "Some summary under N chars",
    "contents": "Some really, really long string"
}'
```

Response: 403

```json
{
  "error": "not authorized to update this resource"
}
```

## Read All Blogs

### Anonymous User Can See all Published Blogs

```sh
curl --location 'http://localhost:8080/api/v1/posts'
```

Response: 200

```json
{
  "posts": [
    {
      "id": "0197aae1-8c00-7e25-80e6-81fc72c7d5f8",
      "status": "PUBLISHED",
      "title": "My riveting blog post with a new title",
      "summary": "Some summary under N chars",
      "contents": "Some really, really long string",
      "author_id": "0197aad2-96f6-7376-8945-18690000ed92",
      "created_ts": "2025-06-25T23:16:37-07:00",
      "published_ts": "",
      "updated_ts": "2025-06-25T23:24:29-07:00"
    }
  ]
}
```

## Delete Blog

### Owner Can Delete

```sh
curl --location --request DELETE 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4'
```

Response: 200

```json
{
  "post_id": "0197aae1-8c00-7e25-80e6-81fc72c7d5f8"
}
```

### Other Users Cannot

```sh
curl --location --request DELETE 'http://localhost:8080/api/v1/posts/0197aae1-8c00-7e25-80e6-81fc72c7d5f8' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3YS04OGQxLWFiMjUzOGJmYzM3YSJ9.mlLNE0sv0YwiKrMYLfy5WmkW9KadVvsjBM-x4WNE7M4'
```

Response: 403

```json
{
  "error": "not authorized to update this resource"
}
```

## Admin Deletes Blog

### Non Admins Do Not Have Access

Even the owner of the post cannot delete via this endpoint

```sh
curl --location --request DELETE 'http://localhost:8080/api/v1/admin/posts/0197aaed-4a35-74da-8574-4165524ac934' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOiIwMTk3YWFkMi05NmY2LTczNzYtODk0NS0xODY5MDAwMGVkOTIifQ.GFwkD_z5V1uQN0OLk_0G_TMBwnfqsuYHq64R7BNsop4'
```

Response: 403

```json
{
  "error": "user is not authorized to perform this action"
}
```

### Admin Has Access

```sh
curl --location --request DELETE 'http://localhost:8080/api/v1/admin/posts/0197aaed-4a35-74da-8574-4165524ac934' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6IjAxOTdhYWQyLTk2ZjYtNzM3YS04OGQxLWFiMjUzOGJmYzM3YSJ9.mlLNE0sv0YwiKrMYLfy5WmkW9KadVvsjBM-x4WNE7M4'
```

Response: 200

```json
{
  "post_id": "0197aaed-4a35-74da-8574-4165524ac934"
}
```
