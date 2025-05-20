# Get all chirps
GET http://localhost:8080/api/chirps

# Get chirp by id
GET http://localhost:8080/api/chirps/b6c7a94c-dfa4-452c-96dc-9c6a7b577daa

# Add chirp
POST http://localhost:8080/api/chirps
{
  "body": "Hello from sharath",
  "user_id": "28ef4485-48ce-4799-84ac-4e9262ef0a84"
}
