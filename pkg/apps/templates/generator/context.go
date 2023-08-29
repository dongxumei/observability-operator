package generator

const (
	StatusAppBuilding       int = iota // Start to build an app, should given the appname and version
	StatusComponentBuilding            // Start to build a component in the app.
	StatusFileBuilding                 // Start to build a file in app layer or component layer
	StatusFileBuilded                  // File builded, to start a new session or return parent
	StatusComponentBuilded             // Component builded, to start a new component session
	StatusAppBuilded                   // App builded, to start the content generating.
	StatusAppGenerateFailed            // App content writing, incase failure.
	StatusAppGenerated                 // Generated, next to finished or modify
)
