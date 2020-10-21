## Go api to create index in elastic search and read config from yaml file

- first install go package to create connection with elastic search

  ```
  go get gopkg.in/olivere/elastic.v7 

  go get github.com/gorilla/mux

- use curl to post json content into elasticsearch 
 
``` 
curl --header "Content-Type: application/json"   --request POST   --data '{"Id":"33","Title":"test", "desc":"content","content":"anything3"}'   http://localhost:10000/article

``` 
