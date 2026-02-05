## Requirement

 - Docker version 28 or newer
 - Go 1.25 or newer

## How to run

 - chmod +x setup.sh
 - ./setup.sh
 - docker compose up --build

## How to use

 - make sure server and client already connected and listening like this
 ![alt text](image.png)
 - hit api using curl or postman | curl "http://localhost:8080/trigger?id=client4"
 - check the server-data directory to ensure file downloaded successfully
