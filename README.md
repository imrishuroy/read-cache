# ReadCache

## Overview
ReadCache is an innovative application designed to help users save, organize, and manage articles, links, and content for future reading or sharing within a community platform. It utilizes Go for the backend, Flutter for the frontend, Postgres for the database, and AWS for cloud services. The app is designed to be intuitive and user-friendly, with a clean and minimalistic UI. It is also highly scalable and can be easily extended to support additional features and functionalities.

UI Repo - https://github.com/imrishuroy/read-cache-ui


## Features
- **Save Links:** Users can easily save articles or links from the internet to read later or share with the community.
- **Organize Content:** Categorize and organize saved content based on user-defined tags or categories for efficient retrieval.
- **Set Reading Reminders:** Set reminders or create to-dos to schedule reading specific articles or links on preferred dates.
- **User Interactions:** Engage with a community by sharing, liking, and commenting on saved content.
- **Secure Storage:** Ensure the security and privacy of user data through robust encryption and authentication mechanisms.
- **Cloud Integration:** Utilize AWS services for scalable cloud hosting, ensuring reliability and accessibility.

## Tech Stack
- **Backend:** Go (Golang)
- **Frontend:** Flutter
- **Database:** Postgres
- **Cloud Services:** AWS (Amazon Web Services)

## Setup Instructions
### Backend (Go)
1. Clone the backend repository.
2. Install required dependencies using `go get`.
3. Configure environment variables for database connections and AWS credentials.
4. Run the server using `go run main.go`.

### Frontend (Flutter)
1. Clone the frontend repository.
2. Install Flutter SDK and necessary dependencies.
3. Set backend API endpoints in the Flutter app to communicate with the backend server.
4. Run the app using `flutter run`.

### Database (Postgres)
1. Install Postgres and create a new database for ReadCache.
2. Set up necessary tables and relations as defined in the backend code.
3. Configure database connection strings in the backend server.

### Cloud (AWS)
1. Set up an AWS account and configure necessary services like S3 for storing user content securely.
2. Manage AWS credentials and access policies for secure interactions between the app and AWS services.

## To Pull Postgress Docker Image
    docker pull postgres
    
## To Run Postgress Docker Image
    docker run --name read-cache -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=IWSIWDF2024 -d postgres

## To download golang-migrate
    brew install golang-migrate

## To download sqlc
   brew install sqlc

## To Run sqlc ( with docker )

    docker pull sqlc/sqlc

    docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate  
    ( Tip - Run in CMD )

## New Migration
    migrate create -ext sql -dir db/migration -seq add_tags

## Resources
    https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
    https://github.com/golang-migrate/migrate - For running migration scripts

    Go Gin - https://go.dev/doc/tutorial/web-service-gin

    Go Viper - https://github.com/spf13/viper

## Mock
    Install https://github.com/uber-go/mock
    Add mockgen to path
    install gomock in your project -  go get github.com/golang/mock/gomock
    Run mockgen -package mockdb -destination db/mock/store.go github.com/imrishuroy/read-cache/db/sqlc Store


## Docker

    1. Build image
      docker build -t read-cache:latest .

    2. Run Container
        docker run --name read-cache-api -p 8080:8080 read-cache:latest   

    3. To Run Container in Production
        docker run --name read-cache-api -p 8080:8080 -e GIN_MODE=release read-cache:latest

    3 To overide default network ip of running database container
        docker run --name read-cache-api -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgres://root:IWSIWDF2024@172.17.0.2:5432/read_cache_db?sslmode=disable" read-cache:latest


    4. To list running container
        docker ps    

    4. To remove docker image
        docker rmi <image-name>

    5. To remove docker container
        docker rm <container-name>

    7 To inspect docker container image
        docker container inspect <container-name>   

    8 To see docker network
        docker network ls

    9 To inspect network
        docker inspect <network-name>                


    

DB_SOURCE=postgres://${RDS_USERNAME:db-user}:${RDS_PASSWORD:password}@${RDS_HOSTNAME:localhost}:${RDS_PORT:5432}/${RDS_DB_NAME:read-cache-db}?sslmode=disable

chmod +x start.sh 
make start.sh executable

## SonarQube
https://github.com/remast/service_sonar

CROS
https://github.com/rs/cors
https://stackoverflow.com/questions/29418478/go-gin-framework-cors

https://docs.sqlc.dev/en/latest/howto/named_parameters.html
