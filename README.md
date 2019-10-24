# Incrementor

Service which provides following gRPC API:
- `Auth` - authorize user to use increment. 
    Expecting args:
    1. username 
    2. password 
    Returns JWT token
- `Register` - just add record in incrementor table for new user
    Expecting args:
    1. username 
    2. password

Next methods must send JWT token for authorization
- `GetNumber`, returns current number in increment
- `IncrementNumber`, increment number by increment size in settings
- `SetSettings`,  
    Expecting args:
    1. increment size (must be positive); 
    2. upper boundary; (when number in increment has equal to upper boundary, number is is sets to zero) 


It's a dummy project built with the only purpose to show my go code when prompted

Run:
- docker-compose -f ./deploy/docker-compose.yaml build
- docker-compose -f ./deploy/docker-compose.yaml up -d
- docker-compose -f ./deploy/docker-compose.yaml exec incrementor /app/migrate up

TODO:
- Version & Build
- Unit tests
- Integration tests
- Docs

