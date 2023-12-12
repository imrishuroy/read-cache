# ReadCache

## Overview
ReadCache is an innovative application designed to help users save, organize, and manage articles, links, and content for future reading or sharing within a community platform. It utilizes Go for the backend, Flutter for the frontend, Postgres for the database, and AWS for cloud services. The app is designed to be intuitive and user-friendly, with a clean and minimalistic UI. It is also highly scalable and can be easily extended to support additional features and functionalities.


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

## To Run sqlc ( with docker )

    docker pull sqlc/sqlc

    docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate  
    ( Tip - Run in CMD )

## Resources
    https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
    https://github.com/golang-migrate/migrate - For running migration scripts
    

