Simple DNS Server implemented in Go
===================================

The Domain Name System (DNS) consists of multiple elements: Authoritative
DNS Servers store and provide DNS record information, Recursive DNS servers
(also referred to as caching DNS servers) are the "middlemen" that recursively
look up information on behalf of an end-user. See
[Authoritative vs. Recursive DNS Servers: What's The Difference](authoritative_recursive)
for an overview.

This project provides a subset of the functionality of an **Authoritative
DNS Server** as a study project. If you need a production-grade DNS Server in Go,
have a look at [CoreDNS](coredns) instead.

Featured on [r/golang](reddit).

![Simple DNS Server implemented in Go](https://raw.githubusercontent.com/dlorch/dnsserver/master/dnsserver-go.gif)

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

TODO
----

* Implement more record types (CNAME, MX, TXT, AAAA, ...) according to section 3.2.2. of [RFC1035](rfc1035)
* Implement [DNS Message Compression](message_compression) according to section 4.1.4. of [RFC1035](rfc1035) (thank you [knome](knome) for pointing this out)

Links
-----

* [RFC 1035: Domain Names - Implementation and Specification](rfc1035)
* [DNS Query Message Format](dns_format)
* [Wireshark](wireshark)
* [Structs Instead of Classes - OOP in Go](structs_classes)
* [Rob Pike - 'Concurrency Is Not Parallelism'](pike_concurrency)

[authoritative_recursive]: http://social.dnsmadeeasy.com/blog/authoritative-vs-recursive-dns-servers-whats-the-difference/
[coredns]: https://coredns.io/
[reddit]: https://www.reddit.com/r/golang/comments/c3n7hl/simple_dns_server_implemented_in_go/
[message_compression]: http://www.tcpipguide.com/free/t_DNSNameNotationandMessageCompressionTechnique-2.htm
[knome]: https://www.reddit.com/r/golang/comments/c3n7hl/simple_dns_server_implemented_in_go/erseh68?utm_source=share&utm_medium=web2x
[rfc1035]: https://www.ietf.org/rfc/rfc1035.txt
[dns_format]: http://www.firewall.cx/networking-topics/protocols/domain-name-system-dns/160-protocols-dns-query.html
[wireshark]: https://www.wireshark.org/
[structs_classes]: https://golangbot.com/structs-instead-of-classes/
[pike_concurrency]: https://www.youtube.com/watch?v=cN_DpYBzKso