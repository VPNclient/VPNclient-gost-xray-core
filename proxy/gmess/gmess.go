// Package gmess contains the implementation of GMess protocol and transportation.
//
// GMess contains both inbound and outbound connections. GMess inbound is usually used on servers
// together with 'freedom' to talk to final destination, while GMess outbound is usually used on
// clients with 'socks' for proxying.
//
// GMess is a GOST-based protocol similar to VMess but uses Russian GOST cryptography standards.
package gmess 