#!/bin/sh

echo "Running pre-start script..."

# Run Redis seed script
/usr/local/bin/redis-seed.sh

# Execute the main command (passed as arguments to the script)
exec "$@"
