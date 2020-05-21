# go-boiler-plate
This is a golang boiler plate setup consisting of the follows.
-> Database: MySql
-> Restful apis for User model
-> Authentication via JWT
-> Authorization for public users, protected users and admins.

## Pre-requisite
1. Go installed
2. MySql installed

## To clone the project
1. Create a dir named bacancy: `mkdir bacancy && cd bacancy`
2. Clone the repo inside bacancy repo: `https://github.com/bacancy/go-boiler-plate.git`
3. Create a new database (mysql) named "root". Connect with username and password "root", the host is 127.0.0.1 and port 3306

## To Run The Project
- Run `go get .` in `go-boiler-plate` folder to install dependancies
- To run the project locally run `go run main.go` in the same repo
- To create a build run `go build`
