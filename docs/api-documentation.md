# Easy Storage API Documentation

This document provides a comprehensive guide to the Easy Storage API endpoints for front-end integration.

## Base URL

All API endpoints are relative to the base URL of your deployment:

```
https://your-api-domain.com
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Error Handling

All endpoints return appropriate HTTP status codes:

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required or invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses follow this format:

```json
{
  "error": "Error message description"
}
```

## API Endpoints

### Authentication

#### Register User

Creates a new user account.

- **URL**: `/api/auth/register`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "storage_quota": 10737418240,
      "storage_used": 0
    },
    "access_token": "jwt-token",
    "refresh_token": "refresh-token",
    "expires_in": 86400
  }
  ```

#### Login

Authenticates a user and returns tokens.

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "storage_quota": 10737418240,
      "storage_used": 1048576
    },
    "access_token": "jwt-token",
    "refresh_token": "refresh-token",
    "expires_in": 86400
  }
  ```

#### Refresh Token

Refreshes an expired access token.

- **URL**: `/api/auth/refresh`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "refresh_token": "refresh-token"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "access_token": "new-jwt-token",
    "refresh_token": "new-refresh-token",
    "expires_in": 86400
  }
  ```

#### Get Current User

Retrieves the current authenticated user's information.

- **URL**: `/api/me`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: `200 OK`
  ```json
  {
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "storage_quota": 10737418240,
      "storage_used": 1048576
    },
    "storage": {
      "quota": 10737418240,
      "used": 1048576,
      "available": 10736369664,
      "used_percentage": 0.01
    }
  }
  ```

#### Change Password

Changes the user's password.

- **URL**: `/api/auth/change-password`
- **Method**: `POST`
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "current_password": "current-password",
    "new_password": "new-password"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Password changed successfully"
  }
  ```

### Files

#### Upload File

Uploads a new file.

- **URL**: `/api/files`
- **Method**: `POST`
- **Auth Required**: Yes
- **Content-Type**: `multipart/form-data`
- **Form Parameters**:
  - `file`: The file to upload
  - `folder_id` (optional): ID of the folder to upload to
- **Success Response**: `201 Created`
  ```json
  {
    "id": "file-id",
    "name": "example.pdf",
    "size": 1048576,
    "content_type": "application/pdf",
    "folder_id": "folder-id",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
  ```

#### List Files

Lists all files for the current user.

- **URL**: `/api/files`
- **Method**: `GET`
- **Auth Required**: Yes
- **Query Parameters**:
  - `limit` (optional): Number of files to return per page (default: 20)
  - `offset` (optional): Number of files to skip (default: 0)
  - `sort` (optional): Field to sort by - `name`, `size`, or `created_at` (default: `created_at`)
  - `sort_dir` (optional): Sort direction - `asc` or `desc` (default: `desc`)
- **Success Response**: `200 OK`
  ```json
  {
    "files": [
      {
        "id": "file-id-1",
        "name": "example1.pdf",
        "size": 1048576,
        "content_type": "application/pdf",
        "folder_id": "folder-id",
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
      },
      {
        "id": "file-id-2",
        "name": "example2.jpg",
        "size": 2097152,
        "content_type": "image/jpeg",
        "folder_id": null,
        "created_at": "2023-01-02T12:00:00Z",
        "updated_at": "2023-01-02T12:00:00Z"
      }
    ],
    "total": 2
  }
  ```

#### Download File

Gets a signed URL to download a file.

- **URL**: `/api/files/:id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**:
  - `id`: ID of the file to download
- **Success Response**: `200 OK`
  ```json
  {
    "url": "https://storage-url.com/signed-url",
    "expires_in": 3600,
    "filename": "example.pdf",
    "content_type": "application/pdf",
    "size": 1048576
  }
  ```

#### Delete File

Deletes a file.

- **URL**: `/api/files/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **URL Parameters**:
  - `id`: ID of the file to delete
- **Success Response**: `204 No Content`

### Folders

#### Create Folder

Creates a new folder.

- **URL**: `/api/folders`
- **Method**: `POST`
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "name": "My Folder",
    "parent_id": "parent-folder-id" // Optional
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": "folder-id",
    "name": "My Folder",
    "parent_id": "parent-folder-id",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
  ```

#### List Folders

Lists folders for the current user.

- **URL**: `/api/folders`
- **Method**: `GET`
- **Auth Required**: Yes
- **Query Parameters**:
  - `parentId` (optional): ID of the parent folder to list contents of
  - `showRootOnly` (optional): If true, only root folders will be returned (default: false)
  - `page` (optional): Page number for pagination (default: 1)
  - `pageSize` (optional): Number of folders per page (default: 10, max: 100)
- **Behavior**:
  - If `parentId` is provided: Returns folders within that specific parent folder
  - If `showRootOnly=true` and no `parentId`: Returns only root-level folders (folders with no parent)
  - If neither `parentId` nor `showRootOnly` is provided: Returns all folders for the user
- **Success Response**: `200 OK`
  ```json
  {
    "folders": [
      {
        "id": "folder-id-1",
        "name": "Folder 1",
        "parent_id": "parent-folder-id",
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
      },
      {
        "id": "folder-id-2",
        "name": "Folder 2",
        "parent_id": "parent-folder-id",
        "created_at": "2023-01-02T12:00:00Z",
        "updated_at": "2023-01-02T12:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 25,
      "total_pages": 3,
      "has_next_page": true,
      "has_prev_page": false
    }
  }
  ```

#### Get Folder Contents

Retrieves all files and folders within a specific folder.

- **URL**: `/api/folders/:folder_id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**:
  - `folder_id`: ID of the folder to get contents of
- **Success Response**: `200 OK`
  ```json
  {
    "folder_id": "folder-id",
    "contents": [
      {
        "id": "folder-id-1",
        "name": "Subfolder",
        "parent_id": "folder-id",
        "type": "folder",
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
      },
      {
        "id": "file-id-1",
        "name": "document.pdf",
        "size": 1048576,
        "content_type": "application/pdf",
        "type": "file",
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
      }
    ],
    "total": 2
  }
  ```

#### Delete Folder

Deletes a folder and all its contents.

- **URL**: `/api/folders/:folder_id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **URL Parameters**:
  - `folder_id`: ID of the folder to delete
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Folder deleted successfully"
  }
  ```

### Shares

#### Create Share

Creates a new share for a file or folder.

- **URL**: `/api/shares`
- **Method**: `POST`
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "resource_id": "file-or-folder-id",
    "resource_type": "file", // or "folder"
    "share_type": "LINK", // or "USER"
    "permission": "READ", // or "WRITE"
    "recipient_id": "user-id", // Required for USER shares
    "password": "optional-password",
    "expires_at": "2023-12-31T23:59:59Z" // Optional expiration date
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": "share-id",
    "resource_id": "file-or-folder-id",
    "resource_type": "file",
    "share_type": "LINK",
    "permission": "READ",
    "recipient_id": null,
    "token": "share-token", // Only for LINK shares
    "has_password": true,
    "expires_at": "2023-12-31T23:59:59Z",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z",
    "access_count": 0,
    "last_access_at": null,
    "is_revoked": false,
    "url": "https://your-domain.com/share/share-token" // Only for LINK shares
  }
  ```

#### List Shares

Lists all shares created by the current user.

- **URL**: `/api/shares`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: `200 OK`
  ```json
  {
    "shares": [
      {
        "id": "share-id-1",
        "resource_id": "file-id-1",
        "resource_type": "file",
        "share_type": "LINK",
        "permission": "READ",
        "token": "share-token-1",
        "has_password": false,
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z",
        "access_count": 5,
        "last_access_at": "2023-01-02T15:30:00Z",
        "is_revoked": false,
        "url": "https://your-domain.com/share/share-token-1"
      },
      {
        "id": "share-id-2",
        "resource_id": "folder-id-1",
        "resource_type": "folder",
        "share_type": "USER",
        "permission": "WRITE",
        "recipient_id": "user-id-2",
        "has_password": false,
        "created_at": "2023-01-02T12:00:00Z",
        "updated_at": "2023-01-02T12:00:00Z",
        "access_count": 0,
        "last_access_at": null,
        "is_revoked": false
      }
    ],
    "total": 2
  }
  ```

#### List Shares With Me

Lists all shares shared with the current user.

- **URL**: `/api/shares/shared-with-me`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: `200 OK`
  ```json
  {
    "shares": [
      {
        "id": "share-id-3",
        "resource_id": "file-id-2",
        "resource_type": "file",
        "share_type": "USER",
        "permission": "READ",
        "recipient_id": "current-user-id",
        "has_password": false,
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z",
        "access_count": 2,
        "last_access_at": "2023-01-02T15:30:00Z",
        "is_revoked": false
      }
    ],
    "total": 1
  }
  ```

#### Get Share

Retrieves a specific share by ID.

- **URL**: `/api/shares/:id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**:
  - `id`: ID of the share to retrieve
- **Success Response**: `200 OK`
  ```json
  {
    "id": "share-id",
    "resource_id": "file-id",
    "resource_type": "file",
    "share_type": "LINK",
    "permission": "READ",
    "token": "share-token",
    "has_password": true,
    "expires_at": "2023-12-31T23:59:59Z",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z",
    "access_count": 10,
    "last_access_at": "2023-01-05T18:45:00Z",
    "is_revoked": false,
    "url": "https://your-domain.com/share/share-token"
  }
  ```

#### Revoke Share

Revokes a share.

- **URL**: `/api/shares/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **URL Parameters**:
  - `id`: ID of the share to revoke
- **Success Response**: `204 No Content`

### Public Share Access

#### Access Shared Resource

Accesses a shared resource using a token.

- **URL**: `/share/:token`
- **Method**: `GET`
- **Auth Required**: No
- **URL Parameters**:
  - `token`: Share token
- **Query Parameters**:
  - `password` (optional): Password for password-protected shares
- **Success Response**: `200 OK`
  ```json
  {
    "id": "share-id",
    "resource_id": "file-id",
    "resource_type": "file",
    "share_type": "LINK",
    "permission": "READ",
    "token": "share-token",
    "has_password": true,
    "expires_at": "2023-12-31T23:59:59Z",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z",
    "access_count": 11,
    "last_access_at": "2023-01-10T14:20:00Z",
    "is_revoked": false,
    "url": "https://your-domain.com/share/share-token"
  }
  ```

#### Download Shared File

Downloads a shared file using a token.

- **URL**: `/share/:token/download`
- **Method**: `GET`
- **Auth Required**: No
- **URL Parameters**:
  - `token`: Share token
- **Query Parameters**:
  - `password` (optional): Password for password-protected shares
- **Success Response**: `200 OK`
  ```json
  {
    "url": "https://storage-url.com/signed-url",
    "expires_in": 3600,
    "filename": "example.pdf",
    "content_type": "application/pdf",
    "size": 1048576
  }
  ```

## Status Codes

The API uses the following status codes:

- `200 OK` - The request was successful
- `201 Created` - A new resource was successfully created
- `204 No Content` - The request was successful but there is no content to return
- `400 Bad Request` - The request was malformed or invalid
- `401 Unauthorized` - Authentication is required or failed
- `403 Forbidden` - The authenticated user does not have permission
- `404 Not Found` - The requested resource was not found
- `500 Internal Server Error` - An error occurred on the server

## Rate Limiting

The API implements rate limiting to prevent abuse. If you exceed the rate limit, you will receive a `429 Too Many Requests` response.

## Pagination

Some endpoints that return lists support pagination using the following query parameters:

- `page`: Page number to return (default: 1)
- `pageSize`: Number of items to return per page (default varies by endpoint)

Paginated responses include a `pagination` object with the following properties:

```json
{
  "current_page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3,
  "has_next_page": true,
  "has_prev_page": false
}
```

This information can be used to build pagination controls in the user interface.

## Sorting

Endpoints that return lists of items often support sorting with the following parameters:

- `sort`: Field to sort by (available fields depend on the endpoint)
- `sort_dir`: Sort direction - `asc` for ascending or `desc` for descending (default: `desc`)

For example, the `/api/files` endpoint supports sorting by `name`, `size`, or `created_at`.

## Versioning

The API is versioned through the URL path. The current version is v1, which is implied in the base URL. 