# What

Wildcard dns resolution for htb machines. Relies on [htb_etc_hosts](https://github.com/fx2301/htb_etc_hosts) being in a sibling directory.

# Example

Get the raw hosts data for HtB:
```
$ git clone https://github.com/fx2301/htb_etc_hosts.git
$ grep bolt.htb htb_etc_hosts/hosts_all.txt
10.10.11.114    bolt.htb             # https://app.hackthebox.com/machines/384
```

Clone this repo and start the DNS server on port 1053:
```
$ git clone https://github.com/fx2301/htb_wildcard_dns.git
$ cd htb_wildcard_dns
$ go run .
Listening at:  :1053
```

```
$ dig bolt.htb @localhost -p 1053 | grep -A 1 SECTION
;; QUESTION SECTION:
;bolt.htb.                      IN      A
--
;; ANSWER SECTION:
bolt.htb.               60      IN      A       10.10.11.114

$ dig fuzzedsubdomain.bolt.htb @localhost -p 1053 | grep -A 1 SECTION
;; QUESTION SECTION:
;fuzzedsubdomain.bolt.htb.      IN      A
--
;; ANSWER SECTION:
bolt.htb.               60      IN      A       10.10.11.114
```

DNS server outputs:
```
Received request from  [::1]:50630
bolt.htb resolved to 10.10.11.114
Received request from  [::1]:53159
fuzzedsubdomain.bolt.htb resolved to 10.10.11.114
```

# Kudos

[dlorch](https://github.com/dlorch) put all the work into making the underpinnings of this repo with [dnsserver](https://github.com/dlorch/dnsserver). Check it out!

# Licensing

Maintained from [dnsserver](https://github.com/dlorch/dnsserver). See [LICENSE](https://github.com/fx2301/htb_wildcard_dns/blob/master/LICENSE).
