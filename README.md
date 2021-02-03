### Tracking Spotify Usage

The aim of this project is a lightweight application to write my spotify history to a database. 

The learning outcome is to understand:
- Use of the Spotify library
- Auth flow
- Writing to a database using Golang

### Deploying 

[Helpful guide](https://medium.com/@migueldoctor/run-mysql-phpmyadmin-locally-in-3-steps-using-docker-74eb735fa1fc) for deploying mysql

Create a [docker volume](https://docs.docker.com/storage/volumes/)

Create a DB
- need to assign key to played_at [(MUL key)](https://www.tutorialspoint.com/create-mysql-column-with-key-mul)

Run mysql docker container:
```bigquery
docker run --name mysql -p 3306:3306 -e "MYSQL_ROOT_PASSWORD=password" -v mysql-volume:/var/lib/mysql -d mysql:latest
```

Connect to the mysql instance:
```bigquery
mysql -h 127.0.0.1 -P 3306 -u root -p
```

### Deploying on a Pi

Follow [this guide](https://www.docker.com/blog/happy-pi-day-docker-raspberry-pi/)