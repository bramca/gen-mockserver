package genmock

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	orderedmapv2 "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/exp/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	keyFile = "key.pem"
	certFile = "cert.pem"

	initServerTemplateHttp = `
const fs = require('fs');
const jsonServer = require('json-server');
const db = require('./%s')
const server = jsonServer.create();
const router = jsonServer.router('./%s');
const middlewares = jsonServer.defaults();
const port = process.env.PORT || %d;
const persistentStorage = process.env.STORE || false;

server.use(jsonServer.bodyParser);
`

	initServerTemplateHttps = `
const https = require("https");
const fs = require('fs');
const jsonServer = require('json-server');
const db = require('./%s')
const server = jsonServer.create();
const router = jsonServer.router('./%s');
const middlewares = jsonServer.defaults();
const port = process.env.PORT || %d;
const persistentStorage = process.env.STORE || false;
const options = {
  key: fs.readFileSync("./%s"),
  cert: fs.readFileSync("./%s")
};

server.use(jsonServer.bodyParser);
`
	rewriterDataTemplate = `
function checkWriteToDb() {
	if (persistentStorage) {
		fs.writeFile('./%s', JSON.stringify(db, undefined, 2), (err) => {
			if (err) throw err;
		});
	}
}

server.use(jsonServer.rewriter({
%s
}));
`
	rewriterTemplate = `	"%s": "%s",`
	serverCallTemplate = `
server.%s('%s', (req, res) => {
	console.log(%s);
	statusCode = %s;
	responseBody = %s;
	%s
	res.status(statusCode).json(responseBody);
});
`
	endServerTemplateHttp = `
server.use(middlewares);
server.use(router);
server.listen(port);
`
	endServerTemplateHttps = `
server.use(middlewares);
server.use(router);
https.createServer(options, server).listen(port);
`

	initPackageJson = `
{
  "name": "%s",
  "main": "%s",
  "scripts": {
    "start": "node %s"
  },
  "dependencies": {
    "json-server": "^0.14.0"
  }
}
`
	initDockerFileTemplate = `
# syntax=docker/dockerfile:1

ARG NODE_VERSION=21.5.0

FROM node:${NODE_VERSION}-alpine

ENV NODE_ENV production

WORKDIR /app

# Copy the rest of the source files into the image.
COPY . .
RUN chmod 777 %s

RUN npm install

# Run the application as a non-root user.
USER node

# Expose the port that the application listens on.
EXPOSE %d

# Run the application.
CMD npm start
`

	initDockerFileTemplateHttps = `
# syntax=docker/dockerfile:1

ARG NODE_VERSION=21.5.0

FROM node:${NODE_VERSION}-alpine

ENV NODE_ENV production

WORKDIR /app

# Copy the rest of the source files into the image.
COPY . .
RUN chmod 777 %s

RUN npm install

RUN apk add openssl

RUN openssl req -x509 -nodes -newkey rsa:2048 -keyout %s -out %s -sha256 -days 365 -subj "/C=NL/ST=Amsterdam/L=Amsterdam/O=Localhost/OU=IT Department/CN=%s" -addext "subjectAltName = DNS:%s"
RUN chmod 644 %s

# Run the application as a non-root user.
USER node

# Expose the port that the application listens on.
EXPOSE %d

# Run the application.
CMD npm start
`

	initDockerComposeTemplate = `
services:
  %s:
    build:
      context: .
    environment:
      NODE_ENV: production
      PORT: %d
      STORE: true
    ports:
      - "%d:%d"
`

	initMakefileTemplate = `
MAKEFLAGS := --no-print-directory --silent

i: install
install: ## install application dependencies
	@npm install

r: run
run: ## run application locally
	@npm start

b: docker.build
docker.build: ## Build the docker container
	@docker build . -t %s

dr: docker.run
docker.run: ## run docker containers
	@docker-compose up -d

ds: docker.stop
docker.stop: ## stop docker containers
	@docker-compose stop
`

	makefileTemplateHttps = `
cr: create.certificate
create.certificate: ## create self signed certificate
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout %s -out %s -sha256 -days 365 \
    		-subj "/C=NL/ST=Amsterdam/L=Amsterdam/O=ING/OU=IT Department/CN=localhost"
	@chmod 644 %s
`
)

type RequestStructure struct {
	Path string
	Method string
	Body string
	DbEntry string
	ResponseCode string
	ResponseBody map[string]any
	RequestParams []string
	RequestBody map[string]any
}

func RandStringBytesRmndr(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    return string(b)
}

func generateExampleData(responseBodyPropertiesSchema *base.Schema) string {
	result := RandStringBytesRmndr(rand.Intn(15))
	enumValues := []string{}
	if responseBodyPropertiesSchema.Enum != nil {
		for _, field := range responseBodyPropertiesSchema.Enum {
			enumValues = append(enumValues, field.Value)
		}
		result = gofakeit.RandomString(enumValues)
	}
	if responseBodyPropertiesSchema.Format != "" {
		switch responseBodyPropertiesSchema.Format {
		case "date-time":
			result = gofakeit.Date().String()
		case "uuid":
			result = gofakeit.UUID()
		case "ip":
			result = gofakeit.IPv4Address()
		case "ip-cidr-block":
			result = fmt.Sprintf("%s/%d", gofakeit.IPv4Address(), gofakeit.IntRange(20, 32))
		case "mac-address":
			result = gofakeit.MacAddress()
		case "address-or-block-or-range":
			result = gofakeit.IPv4Address()
		}
	}

	return result
}

func schemaToPropertyMapV3(schema *base.SchemaProxy, definitions *orderedmapv2.OrderedMap[string, *base.SchemaProxy], responseBody map[string]any, maxRecursion int, recursionDepth int, genExamples bool) map[string]any{
	if recursionDepth > maxRecursion {
		return nil
	}
	responseBodySchema := schema.Schema()
	if schema.IsReference() {
		refSplit := strings.Split(schema.GetReference(), "/")
		responseBodyRef := refSplit[len(refSplit) - 1]
		responseBodyContent := definitions.GetPair(responseBodyRef)
		if responseBodyContent.Value != nil {
			responseBodySchema = responseBodyContent.Value.Schema()
		}
	}
	if responseBodySchema.AllOf != nil {
		for _, schemaField := range responseBodySchema.AllOf {
			responseBody = schemaToPropertyMapV3(schemaField, definitions, responseBody, maxRecursion, recursionDepth, genExamples)
		}
	}
	for responseBodyProperties := responseBodySchema.Properties.First(); responseBodyProperties != nil; responseBodyProperties = responseBodyProperties.Next() {
		responseBodyPropertiesSchema := responseBodyProperties.Value().Schema()
		if responseBodyPropertiesSchema.Type != nil {
			switch responseBodyPropertiesSchema.Type[0] {
			case "string":
				responseBody[responseBodyProperties.Key()] = ""
				if genExamples {
					responseBody[responseBodyProperties.Key()] = generateExampleData(responseBodyPropertiesSchema)
				}
			case "array":
				items := []map[string]any{}
				if responseBodyPropertiesSchema.Items != nil && responseBodyPropertiesSchema.Items.IsA() {
					arrayItemSchema := responseBodyPropertiesSchema.Items.A
					arrayItem := map[string]any{}
					arrayItem = schemaToPropertyMapV3(arrayItemSchema, definitions, arrayItem, maxRecursion, recursionDepth + 1, genExamples)
					if arrayItem != nil {
						items = []map[string]any{arrayItem}
					}
				}
				if len(items) > 0 {
					responseBody[responseBodyProperties.Key()] = items
				} else {
					responseBody[responseBodyProperties.Key()] = []any{}
				}
			case "integer":
				minimum := 0
				maximum := 100
				if responseBodyPropertiesSchema.Minimum != nil {
					minimum = int(*responseBodyPropertiesSchema.Minimum)
				}
				if responseBodyPropertiesSchema.Maximum != nil {
					maximum = int(*responseBodyPropertiesSchema.Maximum)
				}
				responseBody[responseBodyProperties.Key()] = minimum
				if genExamples {
					responseBody[responseBodyProperties.Key()] = gofakeit.IntRange(minimum, maximum)
				}
			case "boolean":
				responseBody[responseBodyProperties.Key()] = false
			case "object":
				responseBody[responseBodyProperties.Key()] = map[string]any{}
				responseBody[responseBodyProperties.Key()] = schemaToPropertyMapV3(responseBodyPropertiesSchema.ParentProxy, definitions, responseBody[responseBodyProperties.Key()].(map[string]any), maxRecursion, recursionDepth, genExamples)
			default:
				responseBody[responseBodyProperties.Key()] = nil
			}
		}
	}

	return responseBody
}

func schemaToPropertyMapV2(schema *base.SchemaProxy, definitions *orderedmap.Map[string, *base.SchemaProxy], responseBody map[string]any, maxRecursion int, recursionDepth int, genExamples bool) map[string]any{
	if recursionDepth > maxRecursion {
		return nil
	}
	responseBodySchema := schema.Schema()
	if schema.IsReference() {
		refSplit := strings.Split(schema.GetReference(), "/")
		responseBodyRef := refSplit[len(refSplit) - 1]
		responseBodyContent := definitions.GetPair(responseBodyRef)
		if responseBodyContent.Value != nil {
			responseBodySchema = responseBodyContent.Value.Schema()
		}
	}
	if responseBodySchema.AllOf != nil {
		for _, schemaField := range responseBodySchema.AllOf {
			responseBody = schemaToPropertyMapV2(schemaField, definitions, responseBody, maxRecursion, recursionDepth, genExamples)
		}
	}
	for responseBodyProperties := responseBodySchema.Properties.First(); responseBodyProperties != nil; responseBodyProperties = responseBodyProperties.Next() {
		responseBodyPropertiesSchema := responseBodyProperties.Value().Schema()
		if responseBodyPropertiesSchema.Type != nil {
			switch responseBodyPropertiesSchema.Type[0] {
			case "string":
				responseBody[responseBodyProperties.Key()] = ""
				if genExamples {
					responseBody[responseBodyProperties.Key()] = generateExampleData(responseBodyPropertiesSchema)
				}
			case "array":
				items := []map[string]any{}
				if responseBodyPropertiesSchema.Items != nil && responseBodyPropertiesSchema.Items.IsA() {
					arrayItemSchema := responseBodyPropertiesSchema.Items.A
					arrayItem := map[string]any{}
					arrayItem = schemaToPropertyMapV2(arrayItemSchema, definitions, arrayItem, maxRecursion, recursionDepth + 1, genExamples)
					if arrayItem != nil {
						items = []map[string]any{arrayItem}
					}
				}
				if len(items) > 0 {
					responseBody[responseBodyProperties.Key()] = items
				} else {
					responseBody[responseBodyProperties.Key()] = []any{}
				}
			case "integer":
				minimum := 0
				maximum := 100
				if responseBodyPropertiesSchema.Minimum != nil {
					minimum = int(*responseBodyPropertiesSchema.Minimum)
				}
				if responseBodyPropertiesSchema.Maximum != nil {
					maximum = int(*responseBodyPropertiesSchema.Maximum)
				}
				responseBody[responseBodyProperties.Key()] = gofakeit.IntRange(minimum, maximum)
			case "boolean":
				responseBody[responseBodyProperties.Key()] = false
			default:
				responseBody[responseBodyProperties.Key()] = nil
			}
		} else {
			responseBody = schemaToPropertyMapV2(responseBodyPropertiesSchema.ParentProxy, definitions, responseBody, maxRecursion, recursionDepth, genExamples)
		}
	}

	return responseBody
}

func SpecV2toRequestStructureMap(specFilename string, maxRecursionDepth int, genExamples bool) map[string]map[string][]RequestStructure {
	api, _ := os.ReadFile(specFilename)

	document, err := libopenapi.NewDocument(api)

	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	docModel, errors := document.BuildV2Model()

	if errors != nil {
		panic(fmt.Sprintf("cannot build doc model: %e", errors))
	}

	featureFileDataStructure := map[string]map[string][]RequestStructure{}
	basePath := docModel.Model.BasePath

	for pathPairs := docModel.Model.Paths.PathItems.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		pathName := pathPairs.Key()
		pathName = fmt.Sprintf("%s%s", basePath, pathName)
		re := regexp.MustCompile(`\{([^}]+)\}`)
		searchParams := re.FindAllStringSubmatch(pathName, -1)
		requestParams := []string{}
		if searchParams != nil {
			for _, param := range searchParams {
				requestParam := strings.ReplaceAll(param[1], "-", "")
				requestParams = append(requestParams, requestParam)
				pathName = strings.ReplaceAll(pathName, param[1], requestParam)
			}
		}
		dbEntry := strings.ReplaceAll(re.ReplaceAllString(pathPairs.Key(), ""), "/", "-")
		dbEntry = strings.ReplaceAll(dbEntry, "--", "-")
		dbEntry = strings.Split(dbEntry, "?")[0]
		if len(dbEntry) > 0 && dbEntry[0] == '-' {
			dbEntry = dbEntry[1:]
		}
		if len(dbEntry) > 0 && dbEntry[len(dbEntry)-1] == '-' {
			dbEntry = dbEntry[:len(dbEntry)-1]
		}
		pathName = re.ReplaceAllString(pathName, ":$1")
		pathItem := pathPairs.Value()
		pathOperations := pathItem.GetOperations()
		for pathOperationPairs := pathOperations.First(); pathOperationPairs != nil; pathOperationPairs = pathOperationPairs.Next() {
			httpMethod := strings.ToLower(pathOperationPairs.Key())
			responseBody := map[string]any{}
			var responseCode string
			for responseCodes := pathOperationPairs.Value().Responses.Codes.First(); responseCodes != nil; responseCodes = responseCodes.Next() {
				responseCodesInt , err := strconv.Atoi(responseCodes.Key())
				if err != nil {
					continue
				}
				if responseCodesInt < 300 {
					responseCode = responseCodes.Key()
					if responseCodes.Value().Schema != nil {
						definitions := docModel.Model.Definitions.Definitions
						responseBody = schemaToPropertyMapV2(responseCodes.Value().Schema, definitions, responseBody, maxRecursionDepth, 0, genExamples)
					}

				}
			}
			if _, ok := featureFileDataStructure[httpMethod]; !ok {
				featureFileDataStructure[httpMethod] = map[string][]RequestStructure{}
			}
			if _, ok := featureFileDataStructure[httpMethod][pathName]; !ok {
				featureFileDataStructure[httpMethod][pathName] = []RequestStructure{}
			}

			req := RequestStructure{
				Path: pathName,
				Method: httpMethod,
				DbEntry: dbEntry,
				ResponseCode: responseCode,
				ResponseBody: responseBody,
				RequestParams: requestParams,
			}
			if len(pathOperationPairs.Value().Parameters) > 0 {
				requestBody := map[string]any{}
				for _, parameter := range pathOperationPairs.Value().Parameters {
					if parameter.In == "body" {
						requestBody = map[string]any{}
						requestBody = schemaToPropertyMapV2(parameter.Schema, docModel.Model.Definitions.Definitions, requestBody, maxRecursionDepth, 0, genExamples)
						req.RequestBody = requestBody
					}
				}
			}
			featureFileDataStructure[httpMethod][pathName] = append(featureFileDataStructure[httpMethod][pathName], req)
			if len(pathOperationPairs.Value().Parameters) > 0 {
				requestBody := map[string]any{}
				for _, parameter := range pathOperationPairs.Value().Parameters {
					if parameter.In == "body" {
						requestBody = map[string]any{}
						requestBody = schemaToPropertyMapV2(parameter.Schema, docModel.Model.Definitions.Definitions, requestBody, maxRecursionDepth, 0, genExamples)
						req.RequestBody = requestBody
					}
				}
				for _, parameter := range pathOperationPairs.Value().Parameters {
					if parameter.In == "query" {
						req = RequestStructure{
							Path: fmt.Sprintf("%s?%s=", pathName, parameter.Name),
							Method: httpMethod,
							DbEntry: dbEntry,
							ResponseCode: responseCode,
							ResponseBody: responseBody,
							RequestParams: requestParams,
							RequestBody: requestBody,
						}
						featureFileDataStructure[httpMethod][pathName] = append(featureFileDataStructure[httpMethod][pathName], req)
					}
				}
			}
		}
	}

	return featureFileDataStructure
}

func SpecV3toRequestStructureMap(specFilename string, maxRecursionDepth int, genExamples bool) map[string]map[string][]RequestStructure {
	api, _ := os.ReadFile(specFilename)

	document, err := libopenapi.NewDocument(api)

	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	docModel, errors := document.BuildV3Model()

	if errors != nil {
		panic(fmt.Sprintf("cannot build doc model: %e", errors))
	}

	featureFileDataStructure := map[string]map[string][]RequestStructure{}

	for pathPairs := docModel.Model.Paths.PathItems.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		pathName := pathPairs.Key()
		re := regexp.MustCompile(`\{([^}]+)\}`)
		searchParams := re.FindAllStringSubmatch(pathName, -1)
		requestParams := []string{}
		if searchParams != nil {
			for _, param := range searchParams {
				requestParam := strings.ReplaceAll(param[1], "-", "")
				requestParams = append(requestParams, requestParam)
				pathName = strings.ReplaceAll(pathName, param[1], requestParam)
			}
		}
		dbEntry := strings.ReplaceAll(re.ReplaceAllString(pathPairs.Key(), ""), "/", "-")
		dbEntry = strings.ReplaceAll(dbEntry, "--", "-")
		dbEntry = strings.Split(dbEntry, "?")[0]
		if len(dbEntry) > 0 && dbEntry[0] == '-' {
			dbEntry = dbEntry[1:]
		}
		if len(dbEntry) > 0 && dbEntry[len(dbEntry)-1] == '-' {
			dbEntry = dbEntry[:len(dbEntry)-1]
		}
		pathName = re.ReplaceAllString(pathName, ":$1")
		pathItem := pathPairs.Value()
		pathOperations := pathItem.GetOperations()
		for pathOperationPairs := pathOperations.First(); pathOperationPairs != nil; pathOperationPairs = pathOperationPairs.Next() {
			httpMethod := strings.ToLower(pathOperationPairs.Key())
			responseBody := map[string]any{}
			var responseCode string
			for responseCodes := pathOperationPairs.Value().Responses.Codes.First(); responseCodes != nil; responseCodes = responseCodes.Next() {
				responseCodesInt , err := strconv.Atoi(responseCodes.Key())
				if err != nil {
					continue
				}
				if responseCodesInt < 300 {
					responseCode = responseCodes.Key()
					if responseCodes.Value().Content.Newest().Value.Schema != nil {
						definitions := docModel.Model.Components.Schemas.OrderedMap
						responseBody = schemaToPropertyMapV3(responseCodes.Value().Content.Newest().Value.Schema, definitions, responseBody, maxRecursionDepth, 0, genExamples)
					}

				}
			}
			if _, ok := featureFileDataStructure[httpMethod]; !ok {
				featureFileDataStructure[httpMethod] = map[string][]RequestStructure{}
			}
			if _, ok := featureFileDataStructure[httpMethod][pathName]; !ok {
				featureFileDataStructure[httpMethod][pathName] = []RequestStructure{}
			}

			req := RequestStructure{
				Path: pathName,
				Method: httpMethod,
				DbEntry: dbEntry,
				ResponseCode: responseCode,
				ResponseBody: responseBody,
				RequestParams: requestParams,
			}


			requestBody := map[string]any{}
			if pathOperationPairs.Value().RequestBody != nil {
				requestBodyContent := pathOperationPairs.Value().RequestBody.Content
				for requestBodyPairs := requestBodyContent.First(); requestBodyPairs != nil; requestBodyPairs = requestBodyPairs.Next() {
					requestBodySchema := requestBodyPairs.Value().Schema
					definitions := docModel.Model.Components.Schemas.OrderedMap
					requestBody = schemaToPropertyMapV3(requestBodySchema, definitions, requestBody, maxRecursionDepth, 0, genExamples)
					req.RequestBody = requestBody
				}
			}
			featureFileDataStructure[httpMethod][pathName] = append(featureFileDataStructure[httpMethod][pathName], req)
			if len(pathOperationPairs.Value().Parameters) > 0 {
				for _, parameter := range pathOperationPairs.Value().Parameters {
					if parameter.In == "query" {
						req = RequestStructure{
							Path: fmt.Sprintf("%s?%s=", pathName, parameter.Name),
							Method: httpMethod,
							DbEntry: dbEntry,
							ResponseCode: responseCode,
							ResponseBody: responseBody,
							RequestParams: requestParams,
							RequestBody: requestBody,
						}
						featureFileDataStructure[httpMethod][pathName] = append(featureFileDataStructure[httpMethod][pathName], req)
					}
				}
			}
		}
	}

	return featureFileDataStructure
}

func GenerateServerFile(scheme string, port int, dbFilename string, serverFilename string, featureFileDataStructure  map[string]map[string][]RequestStructure) {
	featureFile, _ := os.Create(serverFilename)
	defer featureFile.Close()

	featureFileContent := fmt.Sprintf(initServerTemplateHttp, dbFilename, dbFilename, port)
	if scheme == "https" {
		featureFileContent = fmt.Sprintf(initServerTemplateHttps, dbFilename, dbFilename, port, keyFile, certFile)
	}


	var rewriterData []string
	dbEntryMap := map[string][]any{}
	dbEntryCalls := []RequestStructure{}
	dbCallMap := map[string]map[string]bool{}
	for _, calls := range featureFileDataStructure {
		for _, filterPaths := range calls {
			for _, filterPath := range filterPaths {
				if _, ok := dbEntryMap[filterPath.DbEntry]; !ok {
					dbEntryMap[filterPath.DbEntry] = []any{}
				}
				path := fmt.Sprintf("/%s", filterPath.DbEntry)
				if len(filterPath.RequestParams) > 0 {
					path = fmt.Sprintf("%s/:%s", path, strings.Join(filterPath.RequestParams, "/:"))
				}
				rewriterData = append(rewriterData, fmt.Sprintf(rewriterTemplate, filterPath.Path, path))
				addCall := false
				if _, ok := dbCallMap[path]; !ok {
					addCall = true
					dbCallMap[path] = map[string]bool{
						filterPath.Method: true,
					}
				} else if !dbCallMap[path][filterPath.Method] {
					addCall = true
					dbCallMap[path][filterPath.Method] = true
				}
				if addCall {
					dbEntryCalls = append(dbEntryCalls, RequestStructure{
						Path:          path,
						Method:        filterPath.Method,
						Body:	       filterPath.Body,
						DbEntry:       filterPath.DbEntry,
						ResponseBody:  filterPath.ResponseBody,
						ResponseCode:  filterPath.ResponseCode,
						RequestParams: filterPath.RequestParams,
						RequestBody:   filterPath.RequestBody,
					})
				}
			}
		}
	}

	rewriterData[len(rewriterData)-1] = strings.Replace(rewriterData[len(rewriterData)-1], ",", "", 1)
	featureFileContent = fmt.Sprintf("%s%s", featureFileContent, fmt.Sprintf(rewriterDataTemplate, dbFilename, strings.Join(rewriterData, "\n")))

	for _, call := range dbEntryCalls {
		pathline := fmt.Sprintf("/%s", call.DbEntry)
		for _, param := range call.RequestParams {
			pathline = fmt.Sprintf("%s/${req.params.%s}", pathline, param)
		}
		logline := fmt.Sprintf("`%s %s`", strings.ToUpper(call.Method), pathline)
		if call.RequestBody != nil {
			logline = fmt.Sprintf("`%s %s with body ${JSON.stringify(req.body)}`", strings.ToUpper(call.Method), pathline)
		}
		response := "undefined"
		if call.ResponseBody != nil {
			responseJson, _ := json.MarshalIndent(call.ResponseBody, "", "\t\t")
			response = strings.Replace(string(responseJson), "}", "\t}", 1)
		}
		addWriteToDbFunc := ""
		if strings.ToLower(call.Method) != "get" {
			addWriteToDbFunc = "checkWriteToDb();"
		}
		serverCall := fmt.Sprintf(serverCallTemplate, call.Method, call.Path, logline, call.ResponseCode, response, addWriteToDbFunc)
		featureFileContent = fmt.Sprintf("%s%s", featureFileContent, serverCall)
	}
	endServerFile := endServerTemplateHttp
	if scheme == "https" {
		endServerFile = endServerTemplateHttps
	}

	featureFileContent = fmt.Sprintf("%s%s", featureFileContent, endServerFile)

	dbJson, _ := json.MarshalIndent(dbEntryMap, "", "  ")

	dbFile, _ := os.Create(dbFilename)
	defer dbFile.Close()

	_, _ = dbFile.Write(dbJson)

	_, _ = featureFile.Write([]byte(featureFileContent))
}

func GenerateDocker(dbFileName string, serverFile string, port int, scheme string) {
	result := fmt.Sprintf(initDockerFileTemplate, dbFileName, port)
	serverName := strings.Replace(serverFile, ".js", "", 1)
	if scheme == "https" {
		result = fmt.Sprintf(initDockerFileTemplateHttps, dbFileName, keyFile, certFile, serverName, serverName, keyFile, port)
	}
	dockerFile, _ := os.Create("Dockerfile")
	_, _ = dockerFile.Write([]byte(result))
	resultCmp := fmt.Sprintf(initDockerComposeTemplate, serverName, port, port, port)
	dockerCompose, _ := os.Create("compose.yaml")
	_, _ = dockerCompose.Write([]byte(resultCmp))
}

func GeneratePackageJson(serverFile string) {
	result := fmt.Sprintf(initPackageJson, strings.Replace(serverFile, ".js", "", 1), serverFile, serverFile)
	packageJsonFile, _ := os.Create("package.json")
	_, _ = packageJsonFile.Write([]byte(result))
}
