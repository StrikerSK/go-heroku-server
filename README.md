# Custom Server Application
Application developed to gain knowledge of Go programming language and Docker. 
Main focus is to work with database using [GORM](https://gorm.io/) and [PostgreSQL](https://www.postgresql.org/).

### Configuration
Application requires to have the following configurations:
#### Postgres
To set SQL database `DATABASE_URL` environment variable has to be set with following template: `postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable`
#### Redis
To set Redis connection `REDIS_URL` environment variable has to be set with following template: `redis://:@{hostname}}:{port}`
#### Port
Default port number is set to `5000`. To change this property `PORT` environment variable has to be set.