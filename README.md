# Time stamp web server in Go

In this exercise we would like to ask you to create a Go project which solves the
following puzzle:

• You have to implement a small web service, which can store one user-provided unix timestamp (*time.Time) in memory.
• The service must have two endpoints, one for saving the timestamp and an other to fetch it.
• The only allowed content type on the service side is text/plain for both in and egress communications.
• The service must take care of data races (concurrent read-write requests on the timestamp), but mutexes are not allowed. You should find another way to manage the concurrent events.
• In the same process where the service is running, please implement the client side which first stores a timestamp and then reads it back.
• The only output of the application on the standard out (in normal cases) must be the timestamp which it has read in the second step.
• The output of the exercise has to be two source files (main.go and main_test.go). The result must run by executing go run main.go command.
• Test coverage needs to reach at least 2%, maximum allowed coverage is 100%.

Don't panic! We are interested in how far you can go (https://gobyexample.com might
be your friend). 

What we are looking for:
• Has a sensible structure
• Error and edge case handling
• Quality of the test(s)

After all done please describe in a few senteces why you don't recommend to release
the codebase to production.
