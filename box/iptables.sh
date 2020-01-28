#!/bin/sh
iptables -F OUTPUT
ip6tables -F OUTPUT

for i in `seq 0 999`; do
    iptables -A OUTPUT -m owner --uid-owner $i -j ACCEPT
    ip6tables -A OUTPUT -m owner --uid-owner $i -j ACCEPT
done

iptables -A OUTPUT -j DROP
ip6tables -A OUTPUT -j DROP

iptables-save
