Usage

DNBSEQ-T7

```bash
seqbot run \
    -flagptah=/path/to/flag.json \
    -action=wfqtime \
    -dingtoken=token1 \
    -dingtoken=token2

seqbot send \
    -msgfile=message.md \
    -dingtoken=token1

seqbot watch \
    -debug \
    -wfqlog=/path/to/WFQLog \
    -action=wfqtime \
    -action=archive \
    -dingtoken=token1 \
    -dingtoken=token2
```

MGISEQ-2000

```bash
seqbot run \
    -flagpath=/path/to/success.txt \
    -action=uploadtime \
    -dingtoken=token1

seqbot send \
    -msgfile=message.md \
    -dingtoken=token1

# rsyncd
seqbot watch \
    -adapter=inotify \
    -debug \
    -data=/path/to/rsyncd \
    -action=wfqtime \
    -action=uploadtime \
    -action=archive \
    -dingtoken=token1 \
    -dingtoken=token2

# samba
seqbot watch \
    -adapter=scan \
    -interval=300 \
    -debug \
    -data=/path/to/rsyncd \
    -action=uploadtime \
    -action=archive \
    -dingtoken=token1 \
    -dingtoken=token2
```
