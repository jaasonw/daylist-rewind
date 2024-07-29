# daylist-rewind

web application that tracks your daylists as they rotate so they can be brought back as a playlist at any point in time

the application can be viewed at https://daylist-rewind.physicsbirds.dev/login. (pending spotify api access approval)

## demo
https://github.com/user-attachments/assets/b592c001-9e09-4e07-ba96-1b8b6249e462

## stack

- next.js
- shadcn ui
- tailwind
- pocketbase (sqlite)
- golang backend

## Running the App Locally

I'm still working on getting an easily reproducible dev environment but the general steps are

1. create a spotify app from the developer dashboard
2. set `SPOTIFY_ID` and `SPOTIFY_SECRET` in your .env
3. set `REDIRECT_URI` to localhost:3000/api/oauth-callback
4. `docker compose -f docker-compose.dev.yml up` builds and runs the docker containers based on the local directory
5. navigate to localhost:8080/\_/ and create a login then use the `pocketbase_tables.json`
6. create a new admin user and place their set `ADMIN_USER` and `ADMIN_PASSWORD` in .env
7. to develop the frontend/backend independently you can set `BACKEND_URL` to localhost:8080 and run `npm install` then `npm run dev` for the frontend and `cd backend` then `go run main.go` for the backend.
8. For the database, you can set `POCKETBASE_URL` to localhost:8090 and run

```
docker run -d \
  --name daylist-rewind-database-dev \
  --restart unless-stopped \
  -p 8090:8090 \



  -v $(pwd)/pb/pb_data:/pb_data \
  -v $(pwd)/pb/pb_public:/pb_public \
  ghcr.io/muchobien/pocketbase:latest
```
