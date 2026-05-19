# **TL;DR**

A file/text summarizer web app & API written in Go powered by Gemini. Supports
JWT authentication, cursor-based pagination, and file summarization for PDF
and plain text formats.

## Getting Started

*WIP*

## Environment Variables

*WIP*

## Requirements

*WIP*

## **Endpoints**

### **API:** `/api/v1`

| Path                        | Description                                      |
| --------------------------- | ------------------------------------------------ |
| `POST /summarize/file`      | Accepts multipart form file and creates a TL;DR. |
| `POST /summarize/text`      | Accepts text and creates a TL;DR.                |
| `GET /tldrs`                | Returns a paginated list of TL;DRs.              |
| `GET /tldrs/{id}`           | Returns a TL;DR.                                 |
| `PATCH /tldrs/{id}`         | Updates a TL;DR's title from a given ID.         |
| `DELETE /tldrs`             | Batch deletes TL;DRs from a given list of IDs.   |
| `DELETE /tldrs/{id}`        | Deletes a TL;DR from a given ID.                 |
| `GET /me`                   | Returns the authenticated user.                  |
| `PATCH /me/change-username` | Updates the username.                            |
| `PATCH /me/reset-password`  | Updates the password.                            |
| `POST /auth/register`       | Registers a new user.                            |
| `POST /auth/login`          | Authenticates a user.                            |
| `POST /auth/logout`         | Revokes the refresh token.                       |
| `POST /auth/refresh`        | Generates a new refresh token.                   |

### **Admin:** `/admin`

| Path                                | Description                                    |
| ----------------------------------- | ---------------------------------------------- |
| `GET /users`                        | Returns a paginated list of users.             |
| `GET /users/{id}`                   | Returns a user from a given ID.                |
| `PATCH /users/{id}/change-username` | Updates the username of a user.                |
| `PATCH /users/{id}/reset-password`  | Updates the password of a user.                |
| `DELETE /users`                     | Batch deletes users from a given list of IDs.  |
| `DELETE /users/{id}`                | Deletes a user from a given ID.                |
| `DELETE /users/all`                 | Deletes all users (Only for development).      |
| `GET /tldrs`                        | Returns a paginated list of TL;DRs.            |
| `GET /tldrs/{id}`                   | Returns a TL;DR from a given ID.               |
| `PATCH /tldrs/{id}`                 | Updates a TL;DR's title from a given ID.       |
| `DELETE /tldrs`                     | Batch deletes TL;DRs from a given list of IDs. |
| `DELETE /tldrs/{id}`                | Deletes a TL;DR from a given ID.               |
| `DELETE /tldrs/all`                 | Deletes all TL;DRs (Only for development).     |
