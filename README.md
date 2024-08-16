> auth-server is a lightweight auth microservice written in Go using Chi, Bcrypt, & JWT.

#### Structure:

- `/api` - boots up chi router.
- `/user/handler` - registering user routes in chi and related controllers.
- `/user/store` - all db logic concerning anything user related.
- `/auth` - Authentication controllers & JWT logic.
- `/config` - Ensures neccessary values exist on execution.
- `/types` - Handles all global typings.
- `/db` - Mongo connection logic.

#### Build using this technology:

- [Chi](https://github.com/go-chi/chi)
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [JWT](https://jwt.io/)
