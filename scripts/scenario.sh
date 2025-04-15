#! /bin/bash
set -euxo pipefail

ip -n ns1 -6 route add fc00:b:3:: encap seg6 mode encap segs fc00:a:2::,fc00:a:2:0:8000::1234,fc00:b:3:: via fc00:a:2:: dev veth12
ip netns exec ns2 ./build/test --prefix='fc00:a:2:0:8000::/80'

wait
