Release Notes - dpipe - Version v0.5.2.stable
=============================================

### Improvement

### BugFixed


Release Notes - dpipe - Version v0.5.0.rc
=============================================

### New Feature

    . EsOutput periodically report statistics
    . engine permits plugins to register HTTP REST api callback
    . cardinality output will export/reset counters through REST
    . replication from local logs to remote via TCP(compression TODO)
    . init.d script added
    . cardinality output uses gob to checkpoint state between restarts
    . config file validation and visualization[dpconf]
    . alarm when syslog-ng drops messages(ad-hoc parsers)

### Improvement

    . http monitoring with cmd: [ping, stat, plugins]
    . when timestamp is obviously invalid, correct it[als pkg]
    . report msg/s bytes/s
    . ES index patterns add support for ym/ymw/ymd
    . send alarm according to severity instead of accumulated lines num with priority queue
    . accumulate alarm email at night
    . abnormal change severity factor is considered

### BugFixed

    . alslog_input sometimes freeze
      because sortedmap.SortedKeys will dead lock
    . alslog_input.go:260: [/mnt/funplus/logs/fp_rstory/user.0.log]unexpected EOF: fr,1390905013221,{"action":"sell","ide
    . race condition of shared mem, concurrent modification 
      copy on write

### Todo

    . Reload on HUP
    . faster json marshal/unmarshal, currently its 20000ns/op, that's 50K msg/s
    . alarm can compare now with same clock yesterday
    . cardinality of uid doesn't work
    . html alarm email
    . plugins can kill themselves if conf error without polution for others

----
