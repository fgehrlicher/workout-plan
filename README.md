# Workout Plan
⚠️ The Project is still WIP and should not be used for production purposes. ⚠️  
"Workout Plan" is a service with which users can iterate over their own training plans.

## Prerequisites
* docker
* docker-compose
* Bash

## Local Environment
just run ```docker-compose up -d``` to start the plan server, mongodb and swagger server.

## Object Definitions
The plan and exercise definitions can be set in the config. (etc/config.yml)

## Swagger documentation
The Api documentation can be found under http://127.0.0.1:5000/ after the docker-compose stack has been started.

## JWT Token
The plan server does not have an user management and relies on the token for access claims.
The decoded jwt Payload looks something like this:
```json
{
  "access": [
    {
      "Type": "plan",
      "Name": "strengthplan1"
    }
  ],
  "aud": "revelfit.fge.cloud/api/plan",
  "exp": 1591708911,
  "iat": 1591698911,
  "iss": "revelfit.fge.cloud/api/user",
  "nbf": 1591698911,
  "sub": "fgehrlicher"
}
```
"access" indicates which plans are accessible to the user.
The plan server does not validate access claims, so the token issuer should be secure.

To create a new dummy jwt token run:
```
docker-compose run claim-generator
```
