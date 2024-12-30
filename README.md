![code coverage badge](https://github.com/mshortcodes/color_my_practice/actions/workflows/ci.yml/badge.svg)

# Color My Practice 🖌️

## Table of Contents

- [About](#about)
  - [Overview](#overview)
- [Testing the API](#testing-the-api)
  - [Swagger](#swagger)
  - [curl](#curl)
- [API](#api)
  - [Users](#users)
  - [Logs](#logs)
  - [Auth](#auth)
  - [Other](#other)

## About

Color My Practice is a simple app for music students to keep track of their practice time. While I focused on the backend for this project, this is an app I would really like to see fully implemented as I think it would be a fun way for students to visualize their practice time. Building this server was heavily inspired by Boot.dev's Chirpy project, which was one of my favorite projects on Boot.dev. There were many tricky concepts which I wanted to further explore and solidify my understanding of.

Key concepts:

- Building an HTTP server in Go
- Creating a RESTful API for clients to interact with, supporting CRUD operations
- Implementing authentication/authorization
- Working with JWTs, refresh tokens, hashing
- Sending cookies to store tokens and reading them for authentication
- Using SQL to store and retrieve data from a Postgres database
- Deploying to Google Cloud
- Integrating a CI/CD pipeline

Tools:

- PostgreSQL
- Goose
- Sqlc

### Overview

Color My Practice follows a heatmap-style of logging practice time. On the calendar, a single day is clicked up to five times, with each number representing a range of practice times.

For example:

- Very light green -> 0-10 min
- Light green -> 10-20 min
- Green -> 20-30
- Dark green -> 30-60 min
- Very dark green -> 1hr+

This is tracked as "color_depth" in the logs schema, which can be thought of as an enumeration. There is also a "confirmed" column which can only be set to true if the parent has confirmed by entering the password.

## Testing the API

🔓 = unauthenticated endpoint

🔒 = authenticated endpoint

Authentication is handled via cookies.

I've included two ways to interact with the API here: swagger and curl.

### Swagger

Swagger automatically generates interactive documentation when provided a .json or .yml file that defines the API. I chose to implement this for two reasons:

1. Create a visually appealing documentation page
2. Make testing quick and easy

Visit https://colormypractice.com/api/docs to test the API with Swagger.

### Curl

For those who prefer a CLI, I've included example curl requests at the bottom of each endpoint's section. Since this API uses cookies heavily, the example requests will create a file to store/send the cookies in your working directory.

## API

### Users

User resource:

```json
{
  "id": "8f88ab37-133f-411b-bd0f-134c614c390a",
  "created_at": "2024-12-11T16:41:16.609607Z",
  "updated_at": "2024-12-11T16:41:16.609607Z",
  "email": "user@example.com"
}
```

---

#### `POST /api/users` 🔓

Creates a new user.

Request body:

```json
{
  "email": "user@example.com",
  "password": "abc"
}
```

Response body:

```json
{
  "id": "8f88ab37-133f-411b-bd0f-134c614c390a",
  "created_at": "2024-12-11T16:41:16.609607Z",
  "updated_at": "2024-12-11T16:41:16.609607Z",
  "email": "user@example.com"
}
```

---

#### `PUT /api/users` 🔒

Updates a user's email and password.

Request body:

```json
{
  "email": "alice@example.com",
  "password": "abcd"
}
```

Response body:

```json
{
  "id": "8f88ab37-133f-411b-bd0f-134c614c390a",
  "created_at": "2024-12-11T16:41:16.609607Z",
  "updated_at": "2024-12-11T16:58:10.551184Z",
  "email": "alice@example.com"
}
```

---

### Logs

Log resource:

```json
{
  "id": "c8600bd1-6e75-43af-8d7c-bb122c01f541",
  "date": "2024-12-12",
  "color_depth": 5,
  "confirmed": false,
  "user_id": "8f88ab37-133f-411b-bd0f-134c614c390a"
}
```

---

#### `POST /api/logs` 🔒

Creates a practice log for the given day.

- Date must be in YYYY-MM-DD format.
- Color depth must be between 1 and 5.

Request body:

```json
{
  "date": "2024-12-12",
  "color_depth": 5
}
```

Response body:

```json
{
  "id": "c8600bd1-6e75-43af-8d7c-bb122c01f541",
  "date": "2024-12-12",
  "color_depth": 5,
  "confirmed": false,
  "user_id": "8f88ab37-133f-411b-bd0f-134c614c390a"
}
```

---

#### `GET /api/logs` **_INSECURE_** 🔓

Returns an array of all logs in descending order (newest to oldest).

#### `GET /api/logs?user_id=d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c` 🔓

Filters by user when the user ID is provided as a query parameter.

Response body:

```json
[
  {
    "id": "c8600bd1-6e75-43af-8d7c-bb122c01f541",
    "date": "2024-12-12",
    "color_depth": 5,
    "confirmed": true,
    "user_id": "d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c"
  },
  {
    "id": "86a508f5-32a8-41e0-b6c8-660869583efc",
    "date": "2024-12-11",
    "color_depth": 2,
    "confirmed": false,
    "user_id": "d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c"
  }
]
```

If the log ID is provided in the path, then only that log is returned.

#### `GET /api/logs/{logID}` 🔓

```json
{
  "id": "c8600bd1-6e75-43af-8d7c-bb122c01f541",
  "date": "2024-12-12",
  "color_depth": 5,
  "confirmed": false,
  "user_id": "d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c"
}
```

<details>
<summary>curl example</summary>
<code>
curl https://colormypractice.com/api/logs
</code>
</details>

---

#### `PUT /api/logs/confirm` 🔒

Sets the confirmed field to true for all logs given in the array.

Request body:

```json
{
  "log_ids": [
    "3769e508-3dd9-465e-8cda-783d560dfddc",
    "832bed98-a27b-419a-a067-e75d4ba30557"
  ],
  "password": "abcd"
}
```

Response body:

```json
[
  {
    "id": "3769e508-3dd9-465e-8cda-783d560dfddc",
    "date": "2024-12-10",
    "color_depth": 2,
    "confirmed": true,
    "user_id": "d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c"
  },
  {
    "id": "832bed98-a27b-419a-a067-e75d4ba30557",
    "date": "2024-12-09",
    "color_depth": 2,
    "confirmed": true,
    "user_id": "d4eeefe3-0a27-4d72-8c43-32dd02f6cd1c"
  }
]
```

---

#### `DELETE /api/logs/{logID}` 🔒

Deletes a log by ID.

Returns a 204 status code.

---

### Auth

#### `POST /api/login` 🔓

Logs a user in. Sends a JWT and refresh token as cookies.

Request body:

```json
{
  "email": "alice@example.com",
  "password": "abcd"
}
```

Response body:

```json
{
  "id": "8f88ab37-133f-411b-bd0f-134c614c390a",
  "created_at": "2024-12-11T16:41:16.609607Z",
  "updated_at": "2024-12-11T16:58:10.551184Z",
  "email": "alice@example.com"
}
```

---

#### `POST /api/refresh` 🔒

Sends a new JWT after validating the refresh token.

Returns a 204 status code.

---

#### `POST /api/revoke` 🔒

Revokes a refresh token.

Returns a 204 status code.

---

### Other

#### `GET /status` 🔓

Serves a simple status page.

Response body:

Page hits: XXX  
Users: XXX  
Logs: XXX
