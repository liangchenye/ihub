curl -X PUT -F file=@testfile001 localhost:8080/repo/testdata/testfile001

curl localhost:8080/repo/testdata
curl localhost:8080/repo/testdata/testfile001
# we assume testfile002 is added by 'copy'
curl localhost:8080/repo/testdata/testfile002

curl localhost:8080/repometa/testdata
curl localhost:8080/repometa/testdata/testfile001
