#!/bin/sh

# Logic to wait until database 
MAX_TRIES=30
TRIES=0

until pg_isready -h $PGHOST -p $PGPORT -U $PGUSER || [ $TRIES -eq $MAX_TRIES ]; 
    do echo "Verifying if database is up... $TRIES/$MAX_TRIES"
        TRIES=$((TRIES+1))
        sleep 2
done

if [ $TRIES -eq $MAX_TRIES ]; then 
    echo "Database is down after $MAX_TRIES attempts, exiting..."
    exit 1
fi

echo "Database is up."

if [ "$ENV" != "production" ]; then 
    DB_EXISTS=$(psql -h $PGHOST -p $PGPORT -U $PGUSER -tAc "SELECT 1 FROM pg_database WHERE datname='$PGDATABASE'")

    if [ "$DB_EXISTS" != "1" ]; then
        echo "Database $PGDATABASE does not exist. Creating it..."
        createdb -h $PGHOST -p $PGPORT -U $PGUSER $PGDATABASE
    else
        echo "Database $PGDATABASE already exists. Skipping creation."
    fi
fi

# Start the application
if [ "$ENV" = "production" ]; then
  # We only run migrations automatically on prod, so that we can write migrations without them running on every reload in development
  echo "Running migrations to $PGDATABASE db..."
  goose -dir internal/database/migrations postgres "$PG_CONN_STRING" up

  make build && exec ./luso-wiki
else
  exec air
fi
