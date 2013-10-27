package parser

import (
	json "github.com/bitly/go-simplejson"
	"github.com/funkygao/gofmt"
	"time"
)

// Errlog's MongoException parser
type MongodbLogParser struct {
	DbParser
}

// Constructor
func newMongodbLogParser(name string, chAlarm chan<- Alarm, dbFile, createTable, insertSql string) (parser *MongodbLogParser) {
	parser = new(MongodbLogParser)
	parser.init(name, chAlarm, dbFile, createTable, insertSql)

	go parser.collectAlarms()

	return
}

func (this *MongodbLogParser) ParseLine(line string) (area string, ts uint64, data *json.Json) {
	area, ts, data = this.DbParser.ParseLine(line)
	if dryRun {
		return
	}

	cls, err := data.Get("class").String()
	if err != nil || cls != "MongoException" {
		// not a mongodb log
		return
	}

	level, err := data.Get("level").String()
	checkError(err)
	msg, err := data.Get("message").String()
	checkError(err)
	msg = this.normalizeMsg(msg)
	flash, err := data.Get("flash_version_client").String()

	logInfo := extractLogInfo(data)
	this.insert(area, ts, level, msg, flash, logInfo.host)

	return
}

func (this *MongodbLogParser) normalizeMsg(msg string) string {
	r := digitsRegexp.ReplaceAll([]byte(msg), []byte("?"))
	return string(r)
}

func (this *MongodbLogParser) collectAlarms() {
	if dryRun {
		this.chWait <- true
		return
	}

	sleepInterval := time.Duration(this.conf.Int("sleep", 15))
	beepThreshold := this.conf.Int("beep_threshold", 1)
	color := FgCyan + Bright + BgRed

	for {
		time.Sleep(time.Second * sleepInterval)

		this.Lock()
		tsFrom, tsTo, err := this.getCheckpoint("mongo")
		if err != nil {
			this.Unlock()
			continue
		}

		rows := this.query("select count(*) as am, msg from mongo where ts<=? group by msg order by am desc", tsTo)
		parsersLock.Lock()
		this.logCheckpoint(color, tsFrom, tsTo, "MongoException")
		for rows.Next() {
			var msg string
			var amount int64
			err := rows.Scan(&amount, &msg)
			checkError(err)

			if amount >= int64(beepThreshold) {
				this.beep()
			}

			this.colorPrintfLn(color, "%5s %s", gofmt.Comma(amount), msg)
		}
		parsersLock.Unlock()
		rows.Close()

		this.delRecordsBefore("mongo", tsTo)
		this.Unlock()

		if this.stopped {
			this.chWait <- true
			break
		}
	}

}
