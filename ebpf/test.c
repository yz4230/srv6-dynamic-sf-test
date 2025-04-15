// go:build ignore

#include "vmlinux.h"

// clang-format off
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
// clang-format on

#include "log.h"
#include "utils.h"

SEC("lwt_xmit/test")
int test(struct __sk_buff *skb) {
  void *data = (void *)(long)skb->data;
  void *data_end = (void *)(long)skb->data_end;

  struct ipv6hdr *ip6h = data;
  if (data + sizeof(*ip6h) > data_end) {
    bpf_printk("packet truncated");
    return BPF_DROP;
  }

  log("packet arrived: src=%pI6, dst=%pI6", (u64)&ip6h->saddr,
      (u64)&ip6h->daddr);

  if (ip6h->nexthdr != IPPROTO_ROUTING) return BPF_OK;  // 43: Routing Header
  struct ipv6_sr_hdr *sr_hdr = (struct ipv6_sr_hdr *)(ip6h + 1);
  if ((void *)(sr_hdr + 1) > data_end) return BPF_DROP;
  if (sr_hdr->type != 4) return BPF_DROP;  // 4: Segment Routing Header

  u16 func = SID_FUNC(&ip6h->daddr);
  u64 arg = SID_ARG(&ip6h->daddr);

  log("SID: func=%x, arg=%llx", func, arg);

  if (sr_hdr->segments_left == 0) return BPF_DROP;
  sr_hdr->segments_left--;

  struct in6_addr *new_dst_ptr = sr_hdr->segments + sr_hdr->segments_left;
  if ((void *)(new_dst_ptr + 1) > data_end) return BPF_DROP;
  ip6h->daddr = *new_dst_ptr;

  return BPF_LWT_REROUTE;

  return BPF_OK;
}

char _license[] SEC("license") = "GPL";
