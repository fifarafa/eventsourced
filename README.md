<p align="center">
  <img src="logo.png" width="500">
</p>

# Welcome to EventSourced! 🚀

**EventSourced** is your friendly neighborhood library for building event sourcing solutions in Go! If you're looking to add event sourcing patterns to your projects and crave for a tool that blends robustness, efficiency, and simplicity, you've hit the jackpot!

## 📚 Learning Resources

To kickstart your event sourcing journey or deepen your understanding, we've curated a list of essential materials:

- ✅ **[Watch Oskar Dudycz's insightful talk on building an event store](https://youtu.be/gaoZdtQSOTo?si=5fGoIchkE48wZzoX)** - A great introduction to the basics and beyond.
- 🔜 **[Alexey Zimarev's Webinar: You don't need an Event Sourcing framework. Or do you?](https://www.architecture-weekly.com/p/webinar-6-webinar-with-alexey-zimarev)** - A thought-provoking session on the necessity of event sourcing frameworks.
- 🔍 Explore libraries in different programming languages to grasp their strengths and limitations.
- 🧐 Delve into the library from [this GitHub repository](https://github.com/eugene-khyst/postgresql-event-sourcing) and simplify it to its core essence.
- 🏗 Construct a Minimal Viable Product (MVP) using in-memory storage for a straightforward starting point.
- 📘 Dive deep into the nuances of `sql.DB` through [this comprehensive guide](http://go-database-sql.org).

## 🌟 What's Cooking? - Planned Features

**EventSourced** is designed with a vision to cater to a broad range of event sourcing needs, structured across different tiers of functionalities:

### Tier 1 Features - The Essentials
- ✅ **Stream Creation** - Lay the foundation of your event sourcing with stream creation.
- ✅ **Event Tables** - A place for your events to call home.
- ✅ **Appending Events** - We're working on making this process seamless, supporting:
  - Strong consistency for read-after-write peace of mind.
  - Optimistic Concurrency to keep data races at bay.
  - Batch event appending for efficiency.
- 🚧**Event Retrieval** - Fetch events with flexibility, based on stream ID, version, or timestamp.
- **Storage Support** - We love Postgres and DynamoDB, and so does EventSourced! Full integration tests included.

### Tier 2 Features - The Upgrades
- A **FlattenStream** function to snapshot the current state of any stream.
- A **Generic Command Handler** for transforming commands into events, considering the current state.

### Tier 3 Features - The Innovations
- **Projections** for consistent read model updates.
- **Subscriptions** for asynchronous projections.

### Future Ideas 💡
- We're dreaming about supporting Global streams, inspired by Rails Event Store.

## 🎨 Inspiration Corner

Our journey is fueled by the incredible work of others. Here are a few repositories
- https://github.com/EventStore/EventStore-Client-Go
- https://github.com/hallgren/eventsourcing
- https://railseventstore.org/docs/v2/expected_version/

## Questions for myself
- how to handle sql migrations, so it's easy to use for users? how they did it in other libraries?