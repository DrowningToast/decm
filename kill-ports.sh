#!/bin/sh

# Kill processes on ports 3000 (Frontend), 8080 (Backend), 5432 (Database)
echo "Killing processes on ports 3000, 8080, 5432..."

# 3000: Web Frontend
lsof -ti :3000 | xargs kill -9 2>/dev/null

# 8080: Core API Backend
lsof -ti :8080 | xargs kill -9 2>/dev/null

echo "Done."
