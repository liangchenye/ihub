curl -X POST -F file=@ok.txt localhost:8080/oct/lcyok
curl -X POST -F file=@fail.txt localhost:8080/oct/lcyfail

curl localhost:8080/oct/lcyok/image
curl localhost:8080/oct/lcyok/status

curl localhost:8080/oct/lcyfail/image
curl localhost:8080/oct/lcyfail/status
