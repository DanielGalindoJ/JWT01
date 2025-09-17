#!/bin/bash

# Script para ejecutar migraciones

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-auth_db}

DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
MIGRATIONS_PATH="internal/adapters/database/migrations"

echo "🗄️  Ejecutando migraciones..."
echo "Base de datos: $DATABASE_URL"
echo "Migraciones: $MIGRATIONS_PATH"

migrate -path $MIGRATIONS_PATH -database $DATABASE_URL up

if [ $? -eq 0 ]; then
    echo "✅ Migraciones ejecutadas exitosamente"
else
    echo "❌ Error ejecutando migraciones"
    exit 1
fi