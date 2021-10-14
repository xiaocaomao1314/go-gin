kill -9 $(lsof -i:8090 -t)
go build -o main
BUILD_ID=DONTKILLME
nohup ./main &>main.log &
