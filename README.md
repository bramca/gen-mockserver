# 💻 Gen Mockserver

[![Go package](https://github.com/bramca/gen-mockserver/actions/workflows/test.yaml/badge.svg)](https://github.com/bramca/gen-mockserver/actions/workflows/test.yaml)
![GitHub](https://img.shields.io/github/license/bramca/gen-mockserver)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bramca/gen-mockserver)

gen-mockserver is a simple tool for generating a mock server based on an [OpenAPI](https://www.openapis.org/) Spec `v2` or `v3`.

## 🛠️ Installation

`go install github.com/bramca/gen-mockserver/cmd/genmock@latest`
or check out the [releases](https://github.com/bramca/gen-mockserver/releases)

##  🎉 Usage

`go run genmock.go -specfile <path to openapi spec> -specversion <openapi spec version> [-scheme <http (default)|https>] [-port <5000 (default)] [-dbfile <db.json (default)>] [-serverfile <server.js (default)>] [-recursiondepth <0 (default)>] [-exampledata <false (default)>]`

### Options
- `-specfile, -s`
    * path to your openapi specification file
<br><br>
- `-specversion, -v`
    * specify the major version of your spec
    * values: 2, 3
<br><br>
- `-scheme, -c [optional]`
    * specify the scheme that should be used by the mock server
    * values: http (default), https
<br><br>
- `-port, -p [optional]`
    * specify the port that should be used by the mock server
    * values: 5000 (default)
<br><br>
- `-dbfile, -d [optional]`
    * filename for the generated database (use the .json file extension)
    * values: db.json (default)
<br><br>
- `-serverfile, -f [optional]`
    * filename for the generated server (use the .js file extension)
    * values: server.js (default)
<br><br>
- `-recursiondepth, -r [optional]`
    * give the maximum recursion depth to generate the response json (default 0)
    * values: 0 (default)
<br><br>
- `-exampledata, -e [optional]`
    * generate fake example data in the responses
    * values: false (default), true

### Example

`genmock -s openapi.yaml -v 3 -e true -r 1`

This will generate the following files:

```txt
.
├── Dockerfile
├── compose.yaml
├── db.json
├── package.json
└── server.js
```
ref [example folder](./example/)

#### Files
- *Dockerfile*
    * for building a docker container running the server
<br><br>
- *compose.yaml*
    * to run docker compose command creating the container
<br><br>
- *db.json*
    * a database file where you can store mock/example data
<br><br>
- *package.json*
    * npm packages the server depends on
<br><br>
- *server.js*
    * the file that is run when starting the server
    * all the routes are in here with there different operations on them
    * this serves a bit as a template, it will return empty results unless you set the `-e` flag to *true*
    * you can define some logic for every route as you wish in this file

#### Run the server

`npm install ; npm start`
<br>
or
<br>
`docker compose up --build [-d]`
