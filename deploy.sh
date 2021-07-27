kill -9 $(lsof -i:9000 -t)
go build -o main
BUILD_ID=DONTKILLME
nohup ./main &>main.log &
