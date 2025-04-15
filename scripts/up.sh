#! /bin/bash
set -euxo pipefail

# Summary
# network namesaces: ns1, ns2, ns3
# connectivity: ns1:veth12 <-> veth21:ns2:veth23 <-> veth32:ns3
# ip addresses:
#   veth12: fc00:a:1::/32
#   veth21: fc00:a:2::/32
#   veth23: fc00:b:2::/32
#   veth32: fc00:b:3::/32

for ns in ns1 ns2 ns3; do
    ip netns add $ns
    ip netns exec $ns sysctl -w net.ipv6.conf.all.forwarding=1
    ip netns exec $ns sysctl -w net.ipv6.conf.all.seg6_enabled=1
done

ip link add veth12 type veth peer name veth21
ip link add veth23 type veth peer name veth32
ip link set veth12 netns ns1
ip link set veth21 netns ns2
ip link set veth23 netns ns2
ip link set veth32 netns ns3

ip netns exec ns1 sysctl -w net.ipv6.conf.veth12.seg6_enabled=1
ip netns exec ns2 sysctl -w net.ipv6.conf.veth21.seg6_enabled=1
ip netns exec ns2 sysctl -w net.ipv6.conf.veth23.seg6_enabled=1
ip netns exec ns3 sysctl -w net.ipv6.conf.veth32.seg6_enabled=1

ip netns exec ns1 ip -6 addr add fc00:a:1::/32 dev veth12
ip netns exec ns2 ip -6 addr add fc00:a:2::/32 dev veth21
ip netns exec ns2 ip -6 addr add fc00:b:2::/32 dev veth23
ip netns exec ns3 ip -6 addr add fc00:b:3::/32 dev veth32

ip netns exec ns1 ip -6 link set veth12 up
ip netns exec ns2 ip -6 link set veth21 up
ip netns exec ns2 ip -6 link set veth23 up
ip netns exec ns3 ip -6 link set veth32 up
ip netns exec ns1 ip -6 link set lo up
ip netns exec ns2 ip -6 link set lo up
ip netns exec ns3 ip -6 link set lo up

ip netns exec ns1 ip -6 route add default via fc00:a:2:: dev veth12
ip netns exec ns2 ip -6 route add fc00:b::/32 via fc00:a:1:: dev veth21
ip netns exec ns2 ip -6 route add fc00:a::/32 via fc00:b:3:: dev veth23
ip netns exec ns3 ip -6 route add default via fc00:b:2:: dev veth32
