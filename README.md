# Brevity Go

This project aims to be a production-ready template for my future attempts at building an backend server. The codebase should have the following features:

- [ ] OpenAPI design files and toolings to check for backward compatibility.
- [ ] Sensible API design.
- [ ] OpenTelemetry integration: metric, tracing and log.
- [x] Dependency injection for easier testing. May use tools like `wire` to automatically inject dependencies.
- [x] Repository tests that use **real** database.
- [x] Mocking dependencies while testing. Table testing.
- [ ] Sensible way to configure the application for different environment.
- [ ] Sensible database migration mechanism.
- [x] Built into a Docker image with relatively optimized `Dockerfile`.
- [ ] Have a full-fledged local development environment.
- [x] Utilizes CI toolings like Github actions for code quality checking.
- [ ] Have an extendable and maintainable project structure & application architecture.
- [ ] Contains mechanism for keeping track of the application performance.
- [ ] Handling errors in a sensible way.
- [ ] Application should follow best practice for security.
- [x] Should implement graceful shutdown.
- [ ] Should utilize some pattern like caching with Redis.
- [ ] (maybe) E2E tests.