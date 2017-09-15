#!/usr/bin/env bash

DB_NAME=ihub
DB_USER=isula
DB_PASS=isula

echo "USE mysql;\nCREATE DATABASE IF NOT EXISTS ${DB_NAME} DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;\n" | mysql -h127.0.0.1 -uroot -pisula
echo "USE mysql;\nCREATE USER ${DB_USER}@'%' IDENTIFIED BY '${DB_PASS}';\nFLUSH PRIVILEGES;\n" | mysql -h127.0.0.1 -uroot -pisula
echo "USE mysql;\nGRANT ALL PRIVILEGES ON ${DB_NAME}.* TO ${DB_USER}@'%' IDENTIFIED BY '${DB_PASS}';\n" | mysql -h127.0.0.1 -uroot -pisula

echo "finish"
