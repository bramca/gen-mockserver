package main

import (
	"fmt"
	"os"

	genmock "github.com/bramca/gen-mockserver"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	SpecFile         string `short:"s" long:"specfile" description:"[required] path to your openapi specification file" required:"true"`
	SpecMajorVersion int    `short:"v" long:"specversion" choice:"2" choice:"3" description:"[required] specify the major version of your spec" required:"true"`
	Scheme           string `short:"c" long:"scheme" default:"http" choice:"http" choice:"https" description:"[optional] specify the scheme that should be used by the mock server" required:"true"`
	Port             int    `short:"p" long:"port" default:"5000" description:"[optional] specify the port that should be used by the mock server"`
	DbFile           string `short:"d" long:"dbfile" default:"db.json" description:"[optional] filename for the generated database (use the .json file extension)"`
	ServerFile       string `short:"f" long:"serverfile" default:"server.js" description:"[optional] filename for the generated server (use the .js file extension)"`
	RecursionDepth   int    `short:"r" long:"recursiondepth" default:"0" description:"[optional] give the maximum recursion depth to generate the response json (default 0)"`
	GenFakeExamples  bool   `short:"e" long:"exampledata" description:"[optional] generate fake example data in the responses"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with the argument parsing: %v", err)
		os.Exit(1)
	}
	specFile := opts.SpecFile
	specMajorVersion := opts.SpecMajorVersion
	scheme := opts.Scheme
	port := opts.Port
	dbFile := opts.DbFile
	serverFile := opts.ServerFile
	maxRecursionDepth := opts.RecursionDepth
	var featureFileDataStructure map[string]map[string][]genmock.RequestStructure
	if specMajorVersion == 2 {
		featureFileDataStructure, err = genmock.SpecV2toRequestStructureMap(specFile, maxRecursionDepth, opts.GenFakeExamples)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong with parsing the spec file: %v", err)
			os.Exit(1)
		}
	}
	if specMajorVersion == 3 {
		featureFileDataStructure, err = genmock.SpecV3toRequestStructureMap(specFile, maxRecursionDepth, opts.GenFakeExamples)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong with parsing the spec file: %v", err)
			os.Exit(1)
		}
	}

	featureFileContent, err := genmock.GenerateServerFile(scheme, port, dbFile, featureFileDataStructure)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with generating the server file: %e", err)
		os.Exit(1)
	}

	err = genmock.WriteFile(serverFile, []byte(featureFileContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with writing the server file: %e", err)
		os.Exit(1)
	}

	dbFileContent, err := genmock.GenerateDbFile(featureFileDataStructure)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with generating the database file: %e", err)
		os.Exit(1)
	}

	err = genmock.WriteFile(dbFile, []byte(dbFileContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with writing the database file: %e", err)
		os.Exit(1)
	}

	dockerfileContent := genmock.GenerateDockerfile(dbFile, serverFile, port, scheme)
	err = genmock.WriteFile("Dockerfile", []byte(dockerfileContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with writing the docker file: %e", err)
		os.Exit(1)
	}

	dockerComposeContent := genmock.GenerateDockerCompose(serverFile, port)
	err = genmock.WriteFile("compose.yaml", []byte(dockerComposeContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with writing the docker compose file: %e", err)
		os.Exit(1)
	}

	packageJsonContent := genmock.GeneratePackageJson(serverFile)
	err = genmock.WriteFile("package.json", []byte(packageJsonContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with writing the package file: %e", err)
		os.Exit(1)
	}
}
