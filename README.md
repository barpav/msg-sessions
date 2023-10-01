# Sessions (messenger)

This is user authentication and session management component of the [Messenger](https://github.com/barpav/messenger) pet-project.

## Functions

* Starting and ending user sessions.

* Storing session data (Redis).

* Retrieving and updating session info (time started, last activity, IP, client).

* Providing session validation gRPC API for other microservices. 

See microservice [REST API](https://barpav.github.io/msg-api-spec/#/sessions) and [deployment diagram](https://github.com/barpav/messenger#deployment-diagram) for details.