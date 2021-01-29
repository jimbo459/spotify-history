### Tracking Spotify Usage

The aim of this project is a lightweight application to write my spotify history to a database. 

The learning outcome is to understand:
- Use of the Spotify library
- Auth flow
- Writing to a database using Golang

### Deploying 

Create a [docker volume](https://docs.docker.com/storage/volumes/)

Create

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