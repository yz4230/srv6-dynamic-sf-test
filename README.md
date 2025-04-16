# SRv6 動的サービス機能チェーニングのプルーフオブコンセプト

このプロジェクトは、SRv6（Segment Routing v6）を使用した動的サービス機能チェーニング（Dynamic Service Function Chaining）のプルーフオブコンセプト実装です。eBPFとLinuxのネットワーク名前空間を組み合わせて、SRv6のサービスチェーニングを実現します。

## 概要

SRv6は、IPv6ネットワークにおけるセグメントルーティングのプロトコルで、パケットの転送経路を送信元で決定することができます。このプロジェクトでは、eBPFプログラムを使ってSRv6パケットを動的に処理し、サービス機能チェーンを実現します。

## 前提条件

* Linux カーネル 5.4以上（eBPFとSRv6のサポート）
* Go 1.18以上
* Task（タスクランナー）
* bpftool
* clang/LLVM

## セットアップ

### 必要なツールのインストール

```bash
# Ubuntuの場合
sudo apt-get install -y golang clang llvm libelf-dev bpftool
go install github.com/go-task/task/v3/cmd/task@latest
```

### vmlinuxヘッダーの生成

```bash
task dump-vmlinux
```

### ビルド

```bash
task build
```

## 使用方法

### テストシナリオの実行

テストシナリオは次のコマンドで実行できます：

```bash
task scenario
```

この`scenario`タスクは以下の処理を行います：
1. テスト用のネットワーク名前空間を作成（`up.sh`）
2. テストシナリオを実行（`scenario.sh`）
3. 終了時にネットワーク名前空間をクリーンアップ（`down.sh`）

これにより、次の構成のネットワークが作成されます：
- 3つのネットワーク名前空間: ns1, ns2, ns3
- 接続: ns1:veth12 <-> veth21:ns2:veth23 <-> veth32:ns3
- IPアドレス: 
  * veth12: fc00:a:1::/32
  * veth21: fc00:a:2::/32
  * veth23: fc00:b:2::/32
  * veth32: fc00:b:3::/32

### 疎通確認

テストシナリオ実行中に、別のターミナルから以下のコマンドを実行して、ns1からns3への疎通を確認することができます：

```bash
sudo ip netns exec ns1 ping -c 3 fc00:b:3::
```

正常にセットアップされていれば、パケットがns1からns3に届き、pingの応答が返ってくるはずです。これによりSRv6のセグメントルーティングが正常に機能していることを確認できます。

## プロジェクトの構造

```
.
├── build/           - ビルド成果物
├── cmd/             - コマンドラインツール関連コード
├── ebpf/            - eBPFプログラムとヘッダー
│   ├── test.c       - メインのeBPFプログラム
│   ├── utils.h      - ユーティリティ関数
│   ├── log.h        - ロギングヘルパー
│   └── vmlinux.h    - カーネルヘッダー
└── scripts/         - テスト用スクリプト
    ├── up.sh        - ネットワーク名前空間のセットアップ
    ├── scenario.sh  - テストシナリオの実行
    └── down.sh      - ネットワーク名前空間のクリーンアップ
```

## ライセンス

このプロジェクトはGPLライセンスの下で公開されています。
