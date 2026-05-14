# Canary

A simple log aggregator with a dashboard for viewing application logs.

## Todo

- [x] Use serial int64 ids for logs.
- [ ] Use cursor-based pagination instead of Limit/offset.
- [ ] Add identifier in `LogEntry` to be able to identify distinct applications.
- [ ] Allow storing non-json logs. Maybe add a flag to row specifying if the payload is JSON or not.
- [ ] Canary-client: Pipe web-app stdout into the client.
- [ ] App / Client -> Backend communication: Use UDP instead of TCP to potentially improve ingestion speed.
