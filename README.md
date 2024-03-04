<p align="center">
  <img src="logo.png" width="500">
</p>

# eventsourced
`eventsourced` is a project aimed at building an event sourcing library in Go. It is designed to be the default choice for building event sourced systems in Go, providing a robust, efficient, and easy-to-use solution for event sourcing needs.

# Materials
- [X] Watch Oskar Dudycz's talk on building an event store - https://youtu.be/gaoZdtQSOTo?si=5fGoIchkE48wZzoX
- [ ] Watch Alexey Zimarev webinar - You don't need an Event Sourcing framework. Or do you? - https://www.architecture-weekly.com/p/webinar-6-webinar-with-alexey-zimarev
- [ ] Research libraries in other languages, understand their strong points and weaknesses
- [ ] Analyze library from https://github.com/eugene-khyst/postgresql-event-sourcing and descope it to a minimal version
- [ ] Build MVP for MVP (something super simple) - using in-memory storage
- [ ] Learn nitty-gritty of sql.DB from http://go-database-sql.org

# Planned features (in order of priority)
- 1st Tier functionalities
  - CreateStream
  - CreateEventsTable
  - CreateAppendEventFunction
    - It should support strong consistency, so we can read our own writes
    - Putting a simple INSERT, will make us vulnerable to all consistency issues
    - It should support Optimistic Concurrency model
    - In Postgres, using a stored procedure, it makes operations in a single transaction giving us strong consistency
    - It should support appending multiple events at once
  - GetEvents based on the stream ID
    - Support stream version as optional parameter
    - Support timestamp as optional parameter
  - Should support Postgres and DynamoDB as storage
  - Should have integration tests for both storages
- 2nd Tier functionalities 
  - generic FlattenStream function to get the current state of the stream
  - generic command handler, that accepts a command, current state and returns a list of events
- 3rd Tier functionalities
  - Projections - consistent updates of the read model
  - Subscriptions - async projections