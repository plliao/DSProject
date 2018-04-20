# Distributed System Final Project
NYU Tandon 2018 Distributed System Final Project

Instructor: Gustavo Sandoval

## Members
Name | NID | Email
--- | --- | ---
PEI-LUN LIAO | N18410090 | pll273@nyu.edu
I-TING CHEN | N19037964 | itc233@nyu.edu

## Verison

Assignment | Verison
--- | ---:
Part 1: Basic web app | 0.1.0 
Part 2: Separating Front End and Back End | 0.2.0 

## Usage
1. Set up GOPATH

    `export GOPATH=$YOUR_PATH_TO_DSProject/DSProject`

1. Build web application

    `go build main.go`

1. Serve the application

    `./main -port your_port`

    See more instrucitons by
  
    `./main -h`

## API
* Login and Signup

    The entry of the web application is on here with default port number 8080.
    
    `http://hostname:8080/login/`
