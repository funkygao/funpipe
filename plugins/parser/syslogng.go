package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	syslogngDropped = regexp.MustCompile(`dropped=\'program\((.+?)\)=(\d+)\'`)
)

func parseSyslogNgStats(msg string) (match bool, alarm string, severity int) {
	const SYSLOGNG_STATS = "Log statistics; "

	parts := strings.Split(msg, SYSLOGNG_STATS)
	if len(parts) == 2 {
		match = true

		// it is syslog-ng msg in /var/log/messages
		rawStats := parts[1]

		// dropped parsing
		dropped := syslogngDropped.FindAllStringSubmatch(rawStats, 10000)
		for _, d := range dropped {
			num := d[2]
			if num == "0" {
				continue
			}

			// 丢东西啦
			severity = 1000
			alarm = fmt.Sprintf("%s [%s]dropped:%s", alarm, d[1], num)
		}
	} else {
		match = false
	}

	return
}
