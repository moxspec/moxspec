MoxSpec
===

[![CircleCI](https://circleci.com/gh/actapio/moxspec.svg?style=shield&circle-token=1989e1b4d727682e9f80a6fabe13110c282f168c)](https://circleci.com/gh/actapio/moxspec)
[![Maintainability](https://api.codeclimate.com/v1/badges/5492dc939c4a22157f4c/maintainability)](https://codeclimate.com/repos/6015659fc02da6014c007b18/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/5492dc939c4a22157f4c/test_coverage)](https://codeclimate.com/repos/6015659fc02da6014c007b18/test_coverage)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# What is MoxSpec

MoxSpec is a utility to retrieve system hardware information and makes it structured. Minimal footprint and portability by its own hardware interaction engines. 
Information is retrieved and collected directly from standard hardware interfaces such as CPUID, SMBIOS, PCIe configuration register, MSR, and others. Parsing other command-line tools output is not used except for proprietary hardware interfaces such as a hardware RAID because it's slow and unreliable.
It helps performance engineering, hardware review, troubleshooting, asset management, and any other hardware related task.

# Why needed / Story behind MoxSpec
1. Operation at scale with [Open Compute Project (OCP)](https://www.opencompute.org/) / White boxses
When the team adopted Open Compute, instead of benefit out of "part level" replacement, the onsite service team had a challenge troubleshooting / identifying what component needs to be replaced, tool like MoxSpec, allowing SEL decoded to human readable information, was needed to allow onsite service team easily identify what to be replaced such as a DIMM slot.
Democratizing the chance to gain another level of operational scalability to us by adopting open technologies like OCP was crucial and a holistic approach to even a 19” OEM servers was needed.

2. Metrics - Standardization
In a large scale server ops team, life cycle management is one of the key items. In order to get this job done, monitoring of comprehensive hardware metrics such as reads/writes to SSD/NVMe are essential. These are not something IPMI/OOB designed to address. First of all, parsing the output of existing utilities such as smartctl, lspci, .. didn't work out well due to dependency on subcommands, and more importantly those utilities had much more information than needed, which would end up with not a small resource impact when the number of servers become thousands and millions. 
Secondary, in order to support various OSs, the tool needed to be vendor independent and had to take a path of using Linux system programming level commands such as Sys, ioctl et al. This allowed to support many OEM/ODM servers and OSs.
 
![](https://i.imgur.com/xLtWICn.png)


# Example use cases

![](https://i.imgur.com/3GUsEIz.png)

# Prerequisites

- x86_64 or ARMv8 (experimental)
- Intel(Westmere or later) or AMD (Zen or later) processor
- PCI Rev 3.0
- PCI Express Rev 3.0+
- SMBIOS v2.4+
- Linux kernel 2.6.32+

# Development

- Go 1.9 or later

# Installation

```
$ go get github.com/actapio/moxspec/cmd/mox
$ sudo ${GOPATH}/bin/mox
```

# Quick start

```
$ sudo mox show             // standard output
$ sudo mox show -d          // standard output + debug log
$ sudo mox show -j          // output information as a JSON object
```

## Self diagnosis

MoxSpec scans following items for hardware diagnosis and displays `Diag` item as `UNHEALTHY` if a hardware has any errors.
You can diagnose your machine quickly by yourself.

- Processor thermal counter
- Memory ECC error counter
- SMART
- PCIe Advanced Error Report
- RAID controller

`lsdiag` displays diagnosis information.

```
$ sudo lsdiag
+------------+---------+-----------------------------------------+
| category   | stat    | detail                                  |
+------------+---------+-----------------------------------------+
| Processor  | healthy | CPU0 Intel Xeon Gold 6238 2.10GHz       |
| Processor  | healthy | CPU1 Intel Xeon Gold 6238 2.10GHz       |
| Memory     | healthy | Skylake Socket#0 IMC#0 csrow0           |
| Memory     | healthy | Skylake Socket#0 IMC#1 csrow0           |
| Memory     | healthy | Skylake Socket#1 IMC#0 csrow0           |
| Memory     | healthy | Skylake Socket#1 IMC#1 csrow0           |
| NVMe Drive | healthy | Toshiba KXG60ZNV256G TOSHIBA            |
| NVMe Drive | healthy | Toshiba KCM51VUG1T60                    |
| Network    | healthy | Mellanox MT27710 Family [ConnectX-4 Lx] |
| Network    | healthy | Mellanox MT27710 Family [ConnectX-4 Lx] |
+------------+---------+-----------------------------------------+
```

## Detailed RAID information

`lsraid` displays detailed RAID information.
```
$ sudo lsraid
+-----+----------+-----+-----+------+---------+---------+----------------+
| blk | conf     | adp | pos | stat | size    | form    | model          |
+-----+----------+-----+-----+------+---------+---------+----------------+
| sda | RAID 1+0 | 2   | 1:1 | OK   | 900.0GB | SAS HDD | HP EG0900FBVFQ |
| sda | RAID 1+0 | 2   | 1:2 | OK   | 900.0GB | SAS HDD | HP EG0900FBVFQ |
| sda | RAID 1+0 | 2   | 1:3 | OK   | 900.0GB | SAS HDD | HP EG0900FBVFQ |
| sda | RAID 1+0 | 2   | 1:4 | OK   | 900.0GB | SAS HDD | HP EG0900FBVFQ |
+-----+----------+-----+-----+------+---------+---------+----------------+
```

## Serial Number

`lssn` displays serial numbers.
```
$ sudo lssn 
+-----------+-----------------------------------------+--------------------------+----------+--------------------+
| Catagory  | Model                                   | Serial Number            | Location | Spec               |
+-----------+-----------------------------------------+--------------------------+----------+--------------------+
| Processor | Intel Xeon Gold 6238R 2.20GHz           |                          | CPU0     | 28cores, 56threads |
| Processor | Intel Xeon Gold 6238R 2.20GHz           |                          | CPU1     | 28cores, 56threads |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D63                 | DIMM A0  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D64                 | DIMM A1  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D39                 | DIMM A2  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D7C                 | DIMM A3  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D7D                 | DIMM A4  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D5F                 | DIMM A5  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D3E                 | DIMM B0  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D77                 | DIMM B1  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D78                 | DIMM B2  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D7A                 | DIMM B3  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D79                 | DIMM B4  | DDR4-2933 32.0GB   |
| Memory    | Hynix HMA84GR7CJR4N-WM                  | 34F74D7B                 | DIMM B5  | DDR4-2933 32.0GB   |
| Storage   | Toshiba KXG60ZNV256G TOSHIBA            | 496S101TTV7Q             |          | NVMe SSD 256.0GB   |
| Storage   | Toshiba KXD51LN11T92 TOSHIBA            | 90LS1005T51M             |          | NVMe SSD 1.9TB     |
| Storage   | Toshiba KXD51LN11T92 TOSHIBA            | 90LS100ET51M             |          | NVMe SSD 1.9TB     |
| Network   | Mellanox MT27710 Family [ConnectX-4 Lx] |                          |          |                    |
| Network   | Mellanox MT27710 Family [ConnectX-4 Lx] |                          |          |                    |
| System    | Wiwynn SV7220G3                         | P100207834               |          |                    |
| Chassis   |                                         | B9101N010038911000JGAKA1 |          |                    |
| Baseboard | Wistron SV7220G3                        | B5501N010001905000ABJ0B1 |          |                    |
+-----------+-----------------------------------------+--------------------------+----------+--------------------+
```

## Example output

```
System:      MCT TS2290-E7278P-U2 (SN:1947TPS6H0001)
BIOS:        AMI, ver V2.009, release 04/17/2020
Baseboard:   MCT E7278P (SN:MF2D-A00062)
Processor:   2 x Intel Xeon Gold 6238 2.10GHz (44 cores, 88 threads) (2/2 sockets)
               Node
                 package0 node0 (22 cores) (avg 27.3°C) (thermal: safe)
                 package1 node1 (22 cores) (avg 29.7°C) (thermal: safe)
               Cache
                 L1 Data 32KiB/core 8-way
                 L1 Code 32KiB/core 8-way
                 L2 Unified 1024KiB/core 16-way
                 L3 Unified 1408KiB/core 11-way
               TLB
                 L1 Code (4K) 64-entries 8-ways
                 L1 Code (2M/4M) 8-entries fully associative
                 L1 Data (4K) 64-entries 4-ways
                 L1 Data (2M/4M) 32-entries 4-ways
                 L1 Data (1G) 4-entries 4-ways
                 L2 Unified (4K/2M) 1536-entries 12-ways
                 L2 Unified (1G) 16-entries 4-ways
Memory:      Total: 192.0GB
               12 x Samsung DDR4-2933 16.0GB
               4 x empty
             Diag: healthy
Disk:        Intel C620 Series Chipset Family SSATA Controller [AHCI mode] (ahci) (node0)
               Diag: healthy
             Intel C620 Series Chipset Family SATA Controller [AHCI mode] (ahci) (node0)
               Diag: healthy
             Toshiba KXG60ZNV256G TOSHIBA (nvme) (node0) (SN:99CS1001T0KM)
               Link: Gen3 8.0GT/s x4 (max: Gen3 8.0GT/s x4)
               Temp: cur 29°C, warn 78°C, crit 82°C
               Wear: written 1.5TiB, read 1.4TiB (w:51.8%/r:48.2%)
               Firm: AGGA4104
               Diag: healthy
               Namespace: 1
                 nvme0n1 256.0GB (log:512B/phy:512B) (sched:none)
             Toshiba KCM51VUG1T60 (nvme) (node0) (SN:Y9S0A001TVTE)
               Link: Gen3 8.0GT/s x4 (max: Gen3 8.0GT/s x4)
               Temp: cur 32°C, warn 71°C, crit 77°C
               Wear: written 10.0TiB, read 6.6TiB (w:60.3%/r:39.7%)
               Firm: 0107
               Diag: healthy
               Namespace: 1
                 nvme1n1 1.6TB (log:512B/phy:512B) (sched:none)
Network:     Mellanox MT27710 Family [ConnectX-4 Lx] (mlx5_core) (node0)
               Link: Gen3 8.0GT/s x8 (max: Gen3 8.0GT/s x8)
               Intf: enp94s0f0, b8:59:9f:37:9b:94 172.23.140.11/29
               Stat: up, speed 10000, mtu 1500
               Modl: Arista Networks CAB-S-S-25G-2M (SN: ADY1905000DG), SFP Copper DAC (2m)
               Diag: healthy
             Mellanox MT27710 Family [ConnectX-4 Lx] (mlx5_core) (node0)
               Link: Gen3 8.0GT/s x8 (max: Gen3 8.0GT/s x8)
               Intf: enp94s0f1, b8:59:9f:37:9b:95 172.23.140.5/29
               Stat: up, speed 25000, mtu 1500
               Modl: Arista Networks CAB-S-S-25G-2M (SN: ADY1908001WA), SFP Copper DAC (2m)
               Diag: healthy
BMC:         Intf: b8:59:9f:37:9b:96, 10.23.140.5/8
             Firm: 1.17
OS:          CentOS Linux 7 (Core), 3.10.0-1062.18.1.el7.x86_64
Client:      v2.2.12-f113eeb
Hostname:    sample.co.jp
```
```
System:        Wiwynn ST7000G2 N/A (SN:N/A)
BIOS:          AMI, ver BCC17, release 03/23/2020
Baseboard:     Wiwynn BCCIOM001 (SN:B550110700188520007AJ0A1)
Processor:     1 x Intel Xeon D-1581 1.80GHz (16 cores, 32 threads) (1/1 sockets)
                 Node
                   package0 node0 (16 cores) (avg 30.0°C) (thermal: safe)
                 Cache
                   L1 Data 32KiB/core 8-way
                   L1 Code 32KiB/core 8-way
                   L2 Unified 256KiB/core 8-way
                   L3 Unified 1536KiB/core 12-way
                 TLB
                   L1 Code (4K) 64-entries 8-ways
                   L1 Code (2M/4M) 8-entries fully associative
                   L1 Data (4K) 64-entries 4-ways
                   L1 Data (2M/4M) 32-entries 4-ways
                   L1 Data (1G) 4-entries 4-ways
                   L2 Unified (4K/2M) 1536-entries 6-ways
                   L2 Unified (1G) 16-entries 4-ways
Memory:        Total: 64.0GB
                 2 x Hynix Semiconductor DDR4-2667 32.0GB
                 2 x empty
               Diag: healthy
Disk:          Intel 8 Series/C220 Series Chipset Family 6-port SATA Controller 1 [AHCI mode] (ahci) (node0)
                 Diag: healthy
                 Drive: 1
                   sda MTFDDAV256TDL 256.0GB (sd) (log:512B/phy:512B) (sched:deadline) (SN:2002262E1387)
                     Temp: cur 23°C, max 43°C, min 11°C
                     Wear: written 160.9GiB, read 1.3TiB (w:10.5%/r:89.5%)
                     Form: M.2 SSD (SATA 3.3) (cur:6.0Gb/s, max:6.0Gb/s)
                     Firm: M5MU000
                     Diag: healthy
               Broadcom / LSI SAS3008 PCI-Express Fusion-MPT SAS-3 (mpt3sas) (node0)
                 Link: Gen3 8.0GT/s x8 (max: Gen3 8.0GT/s x8)
                 Spec: firm: 13.00.00.00, bios: 13.00.00.00, battery: unknown
                 Diag: healthy
                 Pass-Through Drive: 36
                   [2:0] sdb TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:1] sdc TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:2] sdd TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:3] sde TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:4] sdf TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:5] sdg TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:6] sdh TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:7] sdi TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:8] sdj TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:9] sdk TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:10] sdl TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:11] sdm TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:12] sdn TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:13] sdo TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:14] sdp TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:15] sdq TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:16] sdr TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:17] sds TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:18] sdt TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:19] sdu TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:20] sdv TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:21] sdw TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:22] sdx TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:23] sdy TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:24] sdz TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:25] sdaa TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:26] sdab TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:27] sdac TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:28] sdad TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:29] sdae TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:30] sdaf TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:31] sdag TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:32] sdah TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:33] sdai TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:34] sdaj TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
                   [2:35] sdak TOSHIBA MG07ACA1 14.0TB (sd) (log:512B/phy:4096B) (sched:deadline), Ready (RDY)
Network:       Broadcom BCM57302 NetXtreme-C 10Gb/25Gb Ethernet Controller (bnxt_en) (node0) (SN:00-0a-f7-ff-fe-e6-de-f6)
                 Link: Gen3 8.0GT/s x8 (max: Gen3 8.0GT/s x8)
                 Intf: enp1s0, 00:0a:f7:e6:de:f6 100.81.193.67/26
                 Stat: up, speed 25000, mtu 1500
                 Modl: Arista Networks CAB-S-S-25G-2M (SN: XBM1935100GR), SFP Copper DAC (2m)
                 Diag: healthy
BMC:           Intf: 00:0a:f7:e6:de:f7, 10.81.193.67/8
               Firm: 1.9
OS:            CentOS Linux 7 (Core), 3.10.0-1127.19.1.el7.x86_64
Client:        v2.2.12-f113eeb
Hostname:      sample.co.jp
Last Update:   Sat, 03 Oct 2020 14:15:27 +0900
```
