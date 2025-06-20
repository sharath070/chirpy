# Create users 
POST http://localhost:8080/api/users
{
  "password": "SomeRandomPassword3",
  "email": "example3@email.com"
}

# login 
POST http://localhost:8080/api/login
{
  "password": "SomeRandomPassword3",
  "email": "example3@email.com"
}

# refresh token
POST http://localhost:8080/api/refresh
Authorization: Bearer 3d7b42b63f3518c54a97425a582972d942a1f4983b5bb8ef92d85d2183d9686a
