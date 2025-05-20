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
