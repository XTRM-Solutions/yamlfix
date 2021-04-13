# YAMLFIX

**YAMLFIX supports the OpenAPI
Specification Version 3.0.x**

Postman is a fantastic tool for implementing and testing API
integrations allowing various platforms to take advantage of
the many and various cloud applications and services. These
functions add tremendous flexibility and capabilities to
offerings.

Unfortunately, Postman’s documentation generation from root
OpenAPI Specification files (YAML or JSON) is lacking. OAS
permits both internal and external references to encourage
definition reuse and better organization of code. This
nesting, however, can confuse Postman to the point where
autogenerated schemas are no longer possible. Moreover, for
APIs that communicate primarily through a JSON body
interchange, Postman’s inability to automatically generate
definitions for JSON body parameters renders the automatic
API documentation useless.

YAMLFIX transforms an incoming OAS specification (and
despite its name, YAMLFIX reads both YAML and JSON) in two
ways. Firstly, it steps through POST commands, finds the
basic parameters (that is, parameters that are not
themselves composed of other data) and turns that into two
tables, one for request and one for the response, listing
out the parameters and their definitions. Secondly, it fully
expands the references for each call, allowing Postman to
generate complete example schemas.

## Installing YAMLFIX

- A Go environment must be installed; please
  see [Golang.org](https://www.golang.org) for details.
- A Git command-line client must be installed; please
  see [git-scm.com](https://git-scm.com/) for details
- Open a command console
- Execute `cd %GOPATH\src`
-

Execute `git clone https://<YOUR_GITHUB_USERNAME>:<YOUR_GITHUB_TOKEN>@github.com/XTRM-Solutions/yamlfix.git`

- Execute `cd ./YAMLFIX`
- Execute `go install`  this will compile the program, and
  place the executable into the `%GOROOT%/bin` directory (
  installing Go itself should set both `GOPATH` and `GOROOT`
  environment variables).

## Usage

`yamlfix --infile APISPEC.YAML --outfile APISPEC.JSON`

For more information on options, changing the number of
indent spaces
(or using tabs), debugging, verbose, quiet mode operation,
and option defaults, please run `yamlfix --help`

## Contributing

1. Fork it!
1. Create your feature
   branch: `git checkout -b my-new-feature`
1. Commit your changes: `git commit -am 'Add some feature`
1. Push to the branch: `git push origin my-new-feature`
1. Submit a pull request

### Useful Contribution Areas

* Automated test functionality
* Adding additional output formats (SIMPLEX is overly
  simple, does not map objects to their parameters, etc.)
* Better documentation (always a target)
* Add support for OAS 3.1

## History

| Date | Whom | Company | Comments |
|---|---|---|---|
| 13 April 2021  | Nathan Verrilli  | XTRM | Initial Commit

## Credits

#### Thanks to these two projects

Without the OpenAPI project, this program would be far more
complex.

* github.com/getkin/kin-openapi
* github.com/spf13/pflag

## License

* This code is licensed for use under GPL 3.0
