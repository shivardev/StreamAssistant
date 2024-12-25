code in windows but build for linux as it will take a lot of time 

open cmd and run the following command

set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o alerts

git push in windows, pull in raspberry pi and run the alerts

chmod +x alerts
./alerts