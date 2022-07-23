# Xynnar 

A demo API to help Go begginers that demonstrates the usage of the following:

* GET/POST endpoints using standard library `net/http`
* PostgreSQL
* Caching using Redis
* Generating live API documentation using Swag
* Sort/filter in REST

To checkout the live API with documentation: https://hidden-inlet-47858.herokuapp.com/docs

If the Heroku link doesn't work, email me at muazzam_ali@live.com.

# API

The API has `/api/comments`, `/api/comment`, `/api/films`, and `/api/character` endpoints to read comment, post comment, view films, and view characters.

# Files

`apiHelpers.go` contains the helper methods used by the API. For instance, `sortCharacters()` sorts the character array based on name, height, or gender in ascending or descending order.

`api.go` is where the API endpoints reside for POST and GET calls. Additionally, it has all the data structures for the API.

`main.go` and `utils.go` have utility and initialization functions.
