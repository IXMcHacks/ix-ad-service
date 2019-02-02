# ix-ad-service
A simple demo GO ad-service (web service) for McHacks 2019


### Setup Development Environment
* Setup any directory as your GOPATH
```sh
export GOPATH=$HOME/go
```

* Set the project base path and get inside
```sh
mkdir -p $GOPATH/src/github.com/IXMcHacks
cd $GOPATH/src/github.com/IXMcHacks
```
* Clone the repo
```sh
git clone https://github.com/IXMcHacks/ix-ad-service
cd ix-ad-service
```
* Download Dependencies
```sh
go get github.com/gorilla/schema
```

* Build the project

```sh
go install
```
* Run the project
```sh
 $GOPATH/bin/ix-ad-service
```

### How to install GOLANG
[https://golang.org/doc/install]https://golang.org/doc/install