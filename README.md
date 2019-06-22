Simple DNS Server implemented in Go
===================================

![Simple DNS Server implemented in Go](https://raw.githubusercontent.com/dlorch/dnsserver/master/dnsserver-go.gif)

Featured on [r/golang](https://www.reddit.com/r/golang/comments/c3n7hl/simple_dns_server_implemented_in_go/).

Run
---

```
$ go run dnsserver.go &
Listening at:  :1053

$ dig example.com @localhost -p 1053
Received request from  [::1]:63282

; <<>> DiG 9.10.6 <<>> example.com @localhost -p 1053
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 17060
;; flags: qr; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;example.com.			IN	A

;; ANSWER SECTION:
example.com.		31337	IN	A	3.1.3.7

;; Query time: 0 msec
;; SERVER: ::1#1053(::1)
;; WHEN: Mon Jun 17 17:02:43 CEST 2019
;; MSG SIZE  rcvd: 56
```

Concepts
--------

* Go structs and methods (Go's substitute for object-oriented classes)
* Goroutines ("concurrency is not parallelism" - Rob Pike)
* Go slices (Go's dynamic lists)
* Efficiently writing to and reading from structs using binary.Read() and binary.Write() respectively
* DNS protocol (RFC1035)

Links
-----

* [RFC 1035: Domain Names - Implementation and Specification](https://www.ietf.org/rfc/rfc1035.txt)
* [DNS Query Message Format](http://www.firewall.cx/networking-topics/protocols/domain-name-system-dns/160-protocols-dns-query.html)
* [Wireshark](https://www.wireshark.org/)
* [Structs Instead of Classes - OOP in Go](https://golangbot.com/structs-instead-of-classes/)
* [Rob Pike - 'Concurrency Is Not Parallelism'](https://www.youtube.com/watch?v=cN_DpYBzKso)
