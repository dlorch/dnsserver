Simple DNS Server implemented in Go
===================================

The Domain Name System (DNS) consists of multiple elements: Authoritative
DNS Servers store and provide DNS record information, Recursive DNS servers
(also referred to as caching DNS servers) are the "middlemen" that recursively
look up information on behalf of an end-user. See
[Authoritative vs. Recursive DNS Servers: What's The Difference] for an overview.

This project provides a subset of the functionality of an **Authoritative
DNS Server** as a study project. If you need a production-grade DNS Server in Go,
have a look at [CoreDNS]. For DNS library support, see [Go DNS] or
[package dnsmessage].

Featured on [r/golang] and [go-nuts].

![Simple DNS Server implemented in Go](https://raw.githubusercontent.com/dlorch/dnsserver/master/dnsserver-go.gif)

Run
---

```
$ go run . &
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

* Go structs and methods ([Structs Instead of Classes - OOP in Go])
* Goroutines ([Rob Pike - 'Concurrency Is Not Parallelism'])
* Go slices (Go's dynamic lists)
* Efficiently writing to and reading from structs using binary.Read() and binary.Write() respectively
* DNS protocol ([RFC 1035: Domain Names - Implementation and Specification])

TODO
----

* Implement more record types (CNAME, MX, TXT, AAAA, ...) according to section 3.2.2. of [RFC 1035: Domain Names - Implementation and Specification]
* Implement [DNS Message Compression] according to section 4.1.4. of [RFC 1035: Domain Names - Implementation and Specification] (thank you [knome] for pointing this out)

Links
-----

* [RFC 1035: Domain Names - Implementation and Specification]
* [DNS Query Message Format]
* [Wireshark]
* [Structs Instead of Classes - OOP in Go]
* [Rob Pike - 'Concurrency Is Not Parallelism']

[Authoritative vs. Recursive DNS Servers: What's The Difference]: http://social.dnsmadeeasy.com/blog/authoritative-vs-recursive-dns-servers-whats-the-difference/
[CoreDNS]: https://coredns.io/
[Go DNS]: https://github.com/miekg/dns
[package dnsmessage]: https://godoc.org/golang.org/x/net/dns/dnsmessage
[r/golang]: https://www.reddit.com/r/golang/comments/c3n7hl/simple_dns_server_implemented_in_go/
[go-nuts]: https://groups.google.com/d/msgid/golang-nuts/9d6801ae-5725-4152-83cf-33e63219da70%40googlegroups.com
[DNS Message Compression]: http://www.tcpipguide.com/free/t_DNSNameNotationandMessageCompressionTechnique-2.htm
[knome]: https://www.reddit.com/r/golang/comments/c3n7hl/simple_dns_server_implemented_in_go/erseh68?utm_source=share&utm_medium=web2x
[RFC 1035: Domain Names - Implementation and Specification]: https://www.ietf.org/rfc/rfc1035.txt
[DNS Query Message Format]: http://www.firewall.cx/networking-topics/protocols/domain-name-system-dns/160-protocols-dns-query.html
[Wireshark]: https://www.wireshark.org/
[Structs Instead of Classes - OOP in Go]: https://golangbot.com/structs-instead-of-classes/
[Rob Pike - 'Concurrency Is Not Parallelism']: https://www.youtube.com/watch?v=cN_DpYBzKso
