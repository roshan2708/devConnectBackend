# API Documentation for devConnectBackend

## Overview
This documentation provides detailed information about the API endpoints, database schema, authentication flow, setup instructions, and troubleshooting guide for the devConnectBackend project.

## API Endpoints

### 1. User Authentication
- **POST /api/auth/login**
  - Description: Authenticates a user and retrieves an access token.
- **POST /api/auth/register**
  - Description: Registers a new user.

### 2. User Management
- **GET /api/users**
  - Description: Retrieves a list of all users.
- **GET /api/users/{id}**
  - Description: Retrieves details of a specific user.
- **PUT /api/users/{id}**
  - Description: Updates user information.
- **DELETE /api/users/{id}**
  - Description: Deletes a user.

### 3. Posts
- **GET /api/posts**
  - Description: Retrieves all posts.
- **POST /api/posts**
  - Description: Creates a new post.
- **GET /api/posts/{id}**
  - Description: Retrieves a specific post.
- **PUT /api/posts/{id}**
  - Description: Updates a post.
- **DELETE /api/posts/{id}**
  - Description: Deletes a post.

### 4. Comments
- **GET /api/comments**
  - Description: Retrieves all comments.
- **POST /api/comments**
  - Description: Creates a new comment.
- **GET /api/comments/{id}**
  - Description: Retrieves a specific comment.
- **PUT /api/comments/{id}**
  - Description: Updates a comment.
- **DELETE /api/comments/{id}**
  - Description: Deletes a comment.

### 5. Likes
- **POST /api/likes**
  - Description: Likes a post.
- **DELETE /api/likes/{id}**
  - Description: Removes a like from a post.

### 6. Notifications
- **GET /api/notifications**
  - Description: Retrieves notifications for the logged-in user.

### 7. Search
- **GET /api/search**
  - Description: Searches for users or posts.

## Database Schema
- **Users Table**: Stores user information including username, password, email, etc.
- **Posts Table**: Stores posts with references to the user and content.
- **Comments Table**: Stores comments with references to the post and user.
- **Likes Table**: Stores likes with references to the user and post.

## Authentication Flow
1. A user registers with a username and password.
2. Upon successful registration, the user can log in to receive a token.
3. The token must be included in the header for authentication in protected routes.

## Setup Instructions
1. Clone the repository: `git clone https://github.com/roshan2708/devConnectBackend.git`
2. Navigate to the project directory: `cd devConnectBackend`
3. Install dependencies: `npm install`
4. Start the server: `npm start`
5. Access the API at `http://localhost:5000`

## Troubleshooting Guide
- **Issue**: Unable to connect to the server.
  - **Solution**: Ensure the server is running and you are using the correct port.

- **Issue**: Authentication failures.
  - **Solution**: Verify that your credentials are correct and that you are using the token in requests.

- **Issue**: Missing endpoints.
  - **Solution**: Check the API documentation for the correct routes and methods.

For any additional issues, refer to the project's issue tracker on GitHub or contact the project maintainers.