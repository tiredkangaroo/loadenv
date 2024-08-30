# loadenv

Loadenv provides two options to load from environment files.

## Load
Load loads environment variables from the filepaths specified (defaults to .env).

### Usage:
```golang
package main

import (
  "os"
  "github.com/tiredkangaroo/loadenv"
)

func main() {
  err := loadenv.Load(...filepaths string)
  // handle err here
  os.Getenv("POSTGRES_CONNECTION_URI")
}
```

### Errors:
- reading a file fails
- parsing a line in the file that has bad syntax
- the `setenv()` syscall fails 

## Unmarshal
Unmarshal loads environment variables from the filepaths specified (defaults to .env) as values in a struct.

### Usage:
```golang
package main

import (
  "github.com/tiredkangaroo/loadenv"
)

type EnvironmentVariables struct {
  USERNAME string // defaults to required, will error out if not provided
  SSLMODE  bool `required:"false"` // not required because it is specified in the struct tag
}

func main() {
  var environment EnvironmentVariables
  err := loadenv.Unmarshal(&environment, ...filepaths string)
  // handle err here
  fmt.Println("username:", environment.USERNAME)
}
```

## Expected Environment File Syntax
- All spaces in a line will be trimmed.
- All values are a strings (no quotation marks unless you want them in the string).

### Example
```env
USERNAME=user1
SSLMODE=false
```
