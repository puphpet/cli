package main

import (
	_ "crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		help()
	}

	var action string = os.Args[1]

	switch action {
	case "gen", "generate":
		generate()
		break
	default:
		help()
	}
}

func generate() {
	if len(os.Args) < 3 {
		help_generate()
	}

	var config_file string = os.Args[2]
	var output_file string = ""

	if config_file == "help" {
		help_generate()
	}

	if len(os.Args) < 4 {
		output_file = "puphpet.zip"
	} else {
		output_file = os.Args[3]
	}

	file_contents, err := ioutil.ReadFile(config_file)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	response, err := http.PostForm("http://puphpet.dev/app_dev.php/generate-archive",
		url.Values{"config": {string(file_contents)}})

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("There was an error downloading the zip file.")
		fmt.Println("Check the contents of your YAML file to make sure it is valid.")
		os.Exit(1)
	}

	out, err := os.Create(output_file)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer out.Close()

	if _, err = io.Copy(out, response.Body); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func help() {
	usage := `
Usage:
	puphpet generate [...] => Generates a new zip archive using an existing YAML config file
	puphpet generate help  => Show help for generate command
	puphpet [help]         => This help menu
`

	fmt.Print(usage)
	os.Exit(0)
}

func help_generate() {
	usage := `
Usage:
	puphpet generate config[.yaml] [ output = puphpet.zip ]

Examples:
	puphpet generate config.yaml                => Sends contents of "config.yaml" and receives "puphpet.zip"
	puphpet generate config.yaml downloaded.zip => Sends contents of "config.yaml" and receives "downloaded.zip"
`

	fmt.Print(usage)
	os.Exit(0)
}
