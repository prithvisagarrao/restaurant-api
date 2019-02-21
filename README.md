# Restaurant APIs

### Technology

- Go 1.8
    - Postgres connection details:
        - host: `postgres`
        - port: `5432`
        - dbname: `hellofresh`
        - username: `hellofresh`
        - password: `hellofresh`
- Use the provided `docker-compose.yml` file in the root of this repository. You are free to add more containers to this if you like.


The API conforms to REST practices and provide the following functionality:

- List, create, read, update, and delete Recipes
- Search recipes
- Rate recipes

### Endpoints

Application conforms to the following endpoint structure and returns the HTTP status codes appropriate to each operation. Endpoints specified as protected below require authentication to view. The method of authentication is Json-Web-Tokens(JWT).

##### Recipes

| Name   | Method      | URL                    | Protected |
| ---    | ---         | ---                    | ---       |
| List   | `GET`       | `/recipes`             | ✘         |
| Create | `POST`      | `/recipes`             | ✓         |
| Get    | `GET`       | `/recipes/{id}`        | ✘         |
| Update | `PUT/PATCH` | `/recipes/{id}`        | ✓         |
| Delete | `DELETE`    | `/recipes/{id}`        | ✓         |
| Rate   | `POST`      | `/recipes/{id}/rating` | ✘         |

An endpoint for recipe search functionality is also implemented. The HTTP method and endpoint for this 

### Schema

- **Recipe**
    - Unique ID
    - Name
    - Prep time
    - Difficulty (1-3)
    - Vegetarian (boolean)

Additionally, recipes can be rated many times from 1-5 and a rating is never overwritten.


