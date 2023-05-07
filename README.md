# URL Shortener

## Overview
Create a URL Shortener service that shortens, stores and resolves URLs.

### Example API
- Shorten URL
                    
            curl 'http://localhost:1234/' \
            -H 'Content-Type: text/plain;charset=UTF-8' \
            --data-raw
            'https://someverylongdomainnamehere.com/some/very/very/long/path/here
            ?foo=bar'


<b>Returns "abcde"</b>

- Resolve URL
  
            curl http://localhost:1234/abcde
  

http://localhost:1234/abcde <b>resolves and redirects to</b>
https://someverylongdomainnamehere.com/some/very/very/long/path/here?foo=bar


### Functional requirements
- shortened urls must be short enough (max 5 chars) so they can be easily remembered
-  shortened urls must contain only latin letters and digits
-  if once shortened, shortening the same URL 2+ times must result in the same short URL
string
- the service should be able to store large amounts of URLs (~1 billion)

### Technical requirements
- Implement the solution in the language of your choice.
- Using 3rd party modules and libraries is allowed, but not recommended.
- Provide simple documentation, at least how to set up & run.

### Delivery expectations
- Implement the service following the requirements above.
- Publish the source code into the provided Chaos' internal Git server

#### Bonus points
- Add unit tests.
- Provide detailed documentation for the solution.
The solution will be judged in terms of completeness and quality.
  
## Solution

### Prerequisites:

- Docker with a 'compose' extension - https://docs.docker.com/compose/install/

#### To test the solution execute the following commands from project root:

- Build and enter the container:

            make build

- Run the server:
            
            make run

- Open a new terminal window and navigate to the project,so 
you can enter the container:
  
            make exec
  

- To shorten a url:

            make shorten URL=https://www.verylongurl.com

- To redirect to the long url just add the unique id:

            make redirect URL_ID=aaaaa

- To execute tests:

            make test

- To benchmark the current implementation with 1 million iterations:

            make benchmark

- To benchmark the slow implementation with 100K iterations:

            make benchmark_slow

