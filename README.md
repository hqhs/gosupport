 <!-- https://github.com/madhums/go-gin-mgo-demo -->
 
### Description

### Installation process
 
### CLI DOC

``` shell
$ ./gosupport --help # cli is self documented
```
 
### API DOC 

I'm assuming here what you've created superuser such as "admin@example.com : admin" 

Example 1, everything beyond / and /login is restricted with JWT auth.

```bash
$ http -v --json POST localhost:5050/auth/login email=admin@example.com password=admin
$ http -f GET 
...
HTTP/1.1 200 OK
Content-Length: 185
Content-Type: application/json; charset=utf-8
Date: Wed, 24 Oct 2018 18:31:13 GMT

{
    "nbf": 1540405873,
    "success": "OK",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NDA0MDU4NzMsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.l9mGPO-vDZ59wCTXozsbW2wW1TGu4Xu-CqhvwlUPOkM"
}
$ http -f GET localhost:5050/api/bots "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NDA0MDU4NzMsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.l9mGPO-vDZ59wCTXozsbW2wW1TGu4Xu-CqhvwlUPOkM" "Content-Type: application/json"
```

### Deployment

### Development

#### Tests

Currently, tests doesn't support temprorary databases
