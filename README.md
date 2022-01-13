# Vidz

File Streaming server for your messy computer

Vidz walks your file system in the path provided on the command-line (or defaults to the HOME directory) and serves all the videos on the specified port.

## Building executables

- make linux
- make windows
- make darwin (not tested on Mac)

### Help

./vidz --help

#### Example usage

`./vidz -db=./db.sqlite3 -refresh_db -port=8080`
