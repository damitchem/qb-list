# QB List
"QB List" parser for Infinite Frosthaven Asheron's Call server. Integrates with the already available [AC Quest & Qb List](https://docs.google.com/spreadsheets/d/1RijHs24riB7ww21W7RUfDMdMCi3MqWk8PQCiW-ybryI/edit#gid=737277478).

## Requirements

* Export your full `/qb list` from in-game through your method of choice

Choose one of the following:

If building from source:

* [Go v1.19](https://go.dev/dl/)

If using pre-built:

* Pick from one of the [latest releases](https://github.com/damitchem/qb-list/releases) 

## Recommendations

* Installed [UtilityBelt](https://utilitybelt.gitlab.io/) Decal Plugin for help with exporting the above `/qb list`
    * `/qb list` counts as a `System` message in the UtilityBelt `Chat Log` utility
* Reading the documentation

## Building

* Run `go get` inside of `src`
* Run `go build cmd/qb-list/qb-list.go`

## Getting Started

* Either download one of the [latest releases](https://github.com/damitchem/qb-list/releases) or build from source
* Run the executable via CMD (or any command-line terminal) using appropriate flags defined below
  * Example: qb-list.exe -input=resources/utility-belt-sample.txt -output=C:\\Games -loglevel='Warn'

# Documentation

## Flags

| Flag      | Description                           | Required? |
|-----------|---------------------------------------|-----------|
| -input    | Your input QB file                    | TRUE      |
 | -output   | Output directory, defaults to current | FALSE     |
 | -loglevel | Minimum log level, defaults to Info   | FALSE     |
