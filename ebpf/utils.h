// clang-format off
#include "vmlinux.h" // IWYU pragma: keep
#include <bpf/bpf_helpers.h>
// clang-format on

#ifndef UTILS_H
#define UTILS_H

// IPv6 Routing header
#define IPPROTO_ROUTING 43
#define ntohll(x) (((u64)bpf_ntohl(x)) << 32) + bpf_ntohl(x >> 32)
#define SID_FUNC(x) bpf_ntohs(*((u16 *)(x) + 4))
#define SID_ARG(x) ntohll(*((u64 *)(x) + 1)) & 0x0000FFFFFFFFFFFF

#endif  // UTILS_H
