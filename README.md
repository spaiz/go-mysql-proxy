# go-mysql-proxy
This repository is a result of writing serie of articles `Writing MySQL Proxy` that I'm posting on Medium website.

1 - [Writing MySQL Proxy in GO for self-learning: Part 1 — TCP Proxy](https://medium.com/@alexanderravikovich/quarantine-journey-writing-mysql-proxy-in-go-for-self-learning-part-1-tcp-proxy-39810479b7e9?source=friends_link&sk=9b498aca1d0b239228ab294ba09414bb)
2 - [Writing MySQL Proxy in GO for self-learning: Part 2 — decoding handshake packet](https://medium.com/@alexanderravikovich/writing-mysql-proxy-in-go-for-learning-purposes-part-2-decoding-connection-phase-server-response-7091d87e877e?source=friends_link&sk=c2efb5dfe76e5e061b0679c48e224f2b)

The main goal is to learn the MySQL Protocol by implementing it.

The plan:
- [ ] Implement TCP Proxy as a starting point
- [ ] Implement state machine
- [ ] Implement query/query data buffering
- [ ] Implement plugins

Packets decode/encode todo:
- [ ] Handshake Packet
- [ ] Authorization Packet


go version go1.12.9

To try it, just clone, and run:

```
go run .
```
