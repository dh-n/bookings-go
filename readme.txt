# 🏨 Bookings Go Project - Local Setup Guide

## 1. Local Git Identity Setup
Run these commands inside this specific folder to isolate your commits to your personal account:
```bash
git config --local user.name "Your GitHub Username"
git config --local user.email "your-email@example.com"
```

## 2. Docker Architecture & Controls
This project uses a unified, production-grade `docker-compose.yml` that isolates data inside a persistent named volume (`bookings_postgres_data`) and integrates container health checks.

### Management Workflows via `lazydocker`:
- **Spin Up Database Only:** Highlight `postgres_db` in **Compose Services** and press `u` (Up).
- **Graceful Shutdown:** Disconnect in DBeaver first, then highlight services in `lazydocker` and press `d` (Remove).
- **Hard Database Wipe (Fresh Slate):** If you edit credentials in the compose file and need a clean initialization, run:
  ```bash
  docker compose down
  ```
  ```bash
  docker volume rm bookings_postgres_data
  ```

## 3. Database Connection Parameters (DBeaver)
Use these exact credentials to bridge your local graphical manager to the active container:
- **Host:** `localhost` (or `127.0.0.1`)
- **Port:** `5432`
- **Database:** `bookings`
- **Username:** `postgres`
- **Password:** `mysecretpassword`
- *Note: Ensure "Save password locally" is checked.*

## 4. Database Migrations (Buffalo Fizz)
Once DBeaver passes its connectivity verification check, execute your project migrations natively from your Omarchy terminal:
```bash
soda migrate
```
*(Or use `buffalo db migrate` depending on your course binary path).*
