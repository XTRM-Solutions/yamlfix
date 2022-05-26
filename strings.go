package main

// why this? Because Changing the header level
// is now a one-variable change

const MinIndent int = 1
const MaxIndent int = 32

const HeaderLevel = "3"
const HeaderOpen = "<h" + HeaderLevel + ">"
const HeaderClose = "</h" + HeaderLevel + ">"
const RequestHeader = HeaderOpen + "Request Parameters" + HeaderClose
const ResponseHeader = HeaderOpen + "Response Parameters" + HeaderClose
const SimplexRequestTableHeader = "<table><thead><tr><th>Parameter</th><th>Description</th><th>Required</th></tr></thead>"
const SimplexResponseTableHeader = "<table><thead><tr><th>Parameter</th><th>Description</th></tr></thead>"
const TableBodyOpen = "<tbody>"
const TableBodyAndTableClose = "</tbody></table>"

const KeyAppJson = "application/json"

const TextTrue = "<b>true</b>"
const TextFalse = "<i>false</i>"
