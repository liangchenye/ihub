
## start the testing mysql
```
docker run -p 3306:3306 --net="host" isula/ihub_mysql:0.1
```

## create the test db and test user
```
./mysql/mysql_init.sh
```

## start ihub now
```
docker run --net="host" isula/ihub:0.1
```

