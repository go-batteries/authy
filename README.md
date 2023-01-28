## authy

This library handles access token and refresh token creation. You can revoke the token as well.

Requirements:
- golang
- sqlite3 or postgres (dialects)
- go-migrate cli to write migration files (optional)


The config looks like this:

```go
Config{
    &DatabaseConfig{
        Dialect string
        URL string
    },
    &AppConfig{
        TokenExpiryInSec int64
        RefreshExpiryInSec int64
    }
}
```

### How to use

You need to have a `client_id`, which will generally map to an unique identifier in your data models.
So, for generating `tokens` for `users`, `client_id` will be `user_id`

Presently, `client_id` is string.

##### Initializing and invoking

```go
import (

	"github.com/go-batteries/authy"
	"github.com/go-batteries/authy/pkg/config"
)

var (
    cfg config.Config
    ctx context.Context
    clientID string
)

authService := authy.NewAuthorizer(cfg)
// migrations get run here, where `tokens` table is created

token, err := authService.Create(ctx, clientID)

isValid, err := authService.Authenticate(ctx, token.AccessToken)
_, err = authService.Revoke(ctx, token.AccessToken)

newToken, err := authService.ReAuthenticate(ctx, token.AccessToken, token.RefreshToken)
```

You can also revert the migration as such 

```go
import (
    "database/sql"
    "github.com/go-batteries/authy/database"
)

db, err := sql.Open(/*...*/)
database.RevertMigration(ctx, db)
```


You can check more examples inside `tests/`

TODO:
* [ ] Provide a way to change the table name
