# goauth
basic oauth flows implemented for both client and server

- run the server project which will init the endpoints
- run the client project which will make continous calls to the server's endpoints in both concurrent and sequential form

endpoints:
- /token
- /introspect
- /refresh


basic auth for the authentication is base64(admin:password). This will be moved to use asymmetric key flow.


# TODO:
- enrich token minting
- introspection support
- persistant storage
- middleware tracing
- fallthrough token validation
