# How to get testdata for spc module

## Inquiry Data

```
$ len=`sudo sg_inq /dev/sg1 | grep length | sed -e 's/^\ \+length=\([0-9]\+\).\+$/\1/'`
$ sudo sg_inq /dev/sg1 -r | head -c ${len} | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
00 00 06 12 8b 01 10 02
53 45 41 47 41 54 45 20
53 54 33 30 30 4d 4d 30
30 30 36 20 20 20 20 20
42 30 30 31 53 30 4b 35
4d 5a 5a 34 00 00 00 00
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
00 43 6f 70 79 72 69 67
68 74 20 28 63 29 20 32
30 31 33 20 53 65 61 67
61 74 65 20 41 6c 6c 20
72 69 67 68 74 73 20 72
65 73 65 72 76 65 64 20
```

## Inquiry Data (VPD: Supported VPD pages)

```
$ sudo sg_vpd /dev/sg1 --page=sv -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
00 00 00 10 00 80 83 86
87 88 8a 90 b0 b1 b2 c0
c1 c3 d1 d2
```

## Inquiry Data (VPD: Unit serial number)

```
$ sudo sg_vpd /dev/sg1 --page=sn -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
00 80 00 14 53 30 4b 35
4d 5a 5a 34 30 30 30 30
4b 36 32 38 37 32 55 52
```

## Identify

```
$ sudo sg_sat_identify /dev/sg1 -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
40 04 ff 3f 37 c8 10 00
00 00 00 00 3f 00 00 00
00 00 00 00 20 20 20 20
20 20 20 20 39 31 34 33
33 32 34 37 41 42 31 46
00 00 00 00 00 00 44 20
4d 31 38 55 30 32 69 4d
...
00 00 00 00 00 00 a5 3b
```

## Log EXT

```
# sg_sat_read_gplog is supported in sg3_utils 1.40 or later
$ for line in `sudo sg_sat_read_gplog /dev/sg1 --log=4 --page=1 -HHH | tr -d ' ' | fold -w16`; do
        dump=""
        for word in `echo ${line} | fold -w4`; do
                dump=${dump}`echo ${word} | fold -w2 | tac | tr '\n' ' '`
        done
        echo ${dump}
done

01 00 01 00 00 00 00 00
1c 00 00 00 00 00 00 c0
53 2d 00 00 00 00 00 c0
f0 20 0e 05 00 00 00 c0
3f 47 01 00 00 00 00 c0
35 a1 06 00 00 00 00 c0
...
00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00
```

## Log Sense (Page: Support log pages)

```
$ sudo sg_logs /dev/sg1 --page=0x00 -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
00 00 00 10 00 02 03 05
06 0d 0e 0f 10 15 18 1a
2f 37 38 3e
```

## Log Sense (Page: Error counters (write))

```
$ sudo sg_logs /dev/sg1 --page=0x02 -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
02 00 00 34 00 01 02 04
00 00 00 00 00 02 02 04
00 00 00 00 00 03 02 04
00 00 00 00 00 04 02 04
00 00 00 00 00 05 02 08
00 00 12 21 1c 69 22 00
00 06 02 04 00 00 00 00
```

## Log Sense (Page: Error counters (read))

```
$ sudo sg_logs /dev/sg1 --page=0x03 -r | od -An -tx1 -w8 -v | sed -e 's/^\ //g'
03 00 00 3c 00 00 02 04
a9 c3 3a 5f 00 01 02 04
00 00 00 00 00 02 02 04
00 00 00 00 00 03 02 04
a9 c3 3a 5f 00 04 02 04
00 00 00 00 00 05 02 08
00 00 3e 75 80 f5 92 00
00 06 02 04 00 00 00 00
```

## SMART (Attr/Threshold)

```
$ sudo smartctl -A /dev/sg1 -r ataioctl,2
...
===== [SMART READ ATTRIBUTE VALUES] DATA START (BASE-16) =====
000-015: 10 00 01 2f 00 64 64 00 00 00 00 00 00 00 05 32 |.../.dd........2|
016-031: 00 64 64 00 00 00 00 00 00 00 09 32 00 64 64 53 |.dd........2.ddS|
032-047: 2d 00 00 00 00 00 0c 32 00 64 64 1c 00 00 00 00 |-......2.dd.....|
048-063: 00 00 aa 33 00 64 64 00 00 00 00 00 00 00 ab 32 |...3.dd........2|
...
480-495: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 |................|
496-511: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 a7 |................|
===== [SMART READ ATTRIBUTE VALUES] DATA END (512 Bytes) =====
...
===== [SMART READ ATTRIBUTE THRESHOLDS] DATA START (BASE-16) =====
000-015: 10 00 01 32 00 00 00 00 00 00 00 00 00 00 05 01 |...2............|
016-031: 00 00 00 00 00 00 00 00 00 00 09 00 00 00 00 00 |................|
032-047: 00 00 00 00 00 00 0c 01 00 00 00 00 00 00 00 00 |................|
048-063: 00 00 aa 0a 00 00 00 00 00 00 00 00 00 00 ab 00 |................|
...
480-495: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 |................|
496-511: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 0e |................|
===== [SMART READ ATTRIBUTE THRESHOLDS] DATA END (512 Bytes) =====
...
```
