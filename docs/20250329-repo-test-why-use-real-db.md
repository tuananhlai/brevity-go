# Why use a real database instance for your repository tests?

Outline:
- Traditionally, checking that a repository method work is very annoying. If the database query it use were very simple, then you can try running it directly. Otherwise, you have to wait until you've finished building the service and controller layers before you can start running the query.
- There are multiple ways to test repository methods:
    - Mock out the database. Basically, you can only check that the called database query string matches a certain pattern. Example of this is `sqlmock` package in Go.
    - Use in-memory database. You test with a real database instance, but it's not the same database as your application. Therefore, some methods will fail in tests but works in production and vice versa.
    - Use a real database instance. You test with the same database as your application. This way, the tests you write would be much more reliable. But it comes with some performance drawbacks during test.