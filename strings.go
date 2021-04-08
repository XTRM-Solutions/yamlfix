package main

// why this? Because Changing the header level
// is now a one-variable change
const HeaderLevel = "3"
const HeaderOpen = "<h" + HeaderLevel + ">"
const HeaderClose = "</h" + HeaderLevel + ">"
const RequestHeader = HeaderOpen + "Request Parameters" + HeaderClose
const ResponseHeader = HeaderOpen + "Response Parameters" + HeaderClose
const SimplexTableHeader = "<table><thead><tr><th>Parameter</th><th>Description</th><th>Required</th></tr></thead>"
const TableOpen = "<tbody>"
const TableClose = "</tbody></table>"

const KeyAppJson = "application/json"
