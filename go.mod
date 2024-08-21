module example.com

go 1.22.2

require (
	github.com/gobuffalo/envy v1.10.2
	router v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
)

replace router => ./src/router
