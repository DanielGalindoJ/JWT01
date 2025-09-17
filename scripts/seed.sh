#!/bin/bash

# Script para insertar datos de prueba

echo "🌱 Insertando datos de prueba..."

curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123",
    "first_name": "Test",
    "last_name": "User",
    "username": "testuser"
  }'

echo "\n✅ Datos de prueba insertados"
