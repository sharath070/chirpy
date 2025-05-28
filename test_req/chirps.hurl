# Get all chirps
GET http://localhost:8080/api/chirps

# Get chirp by id
GET http://localhost:8080/api/chirps/b6c7a94c-dfa4-452c-96dc-9c6a7b577daa

# Add chirp
POST http://localhost:8080/api/chirps
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJjMDBiZWU3ZC0zYjJjLTQ4YmItODczZS0yYzcxZmIxY2M4ZTciLCJleHAiOjE3NDgyMTM4MTAsImlhdCI6MTc0ODIxMDIxMH0.LfHSNYyifQRFWYZfof-ZyBoNm7O-flJaWYKK8TyZzPs

{
  "body": "Hellow brow... howw arree you dooingg... can'd you guess by my aksent... I am INDIAN"
}
