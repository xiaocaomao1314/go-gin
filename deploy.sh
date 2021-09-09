kill -9 $(lsof -i:9999 -t)
go build -o main
BUILD_ID=DONTKILLME
nohup ./main &>main.log &
