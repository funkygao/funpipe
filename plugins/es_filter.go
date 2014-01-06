package plugins

import (
	"fmt"
	"github.com/funkygao/als"
	"github.com/funkygao/funpipe/engine"
	conf "github.com/funkygao/jsconf"
)

type converter struct {
	key    string // key name
	action string // action
	cur    string // currency
	rang   []int  // range
}

func (this *converter) load(section *conf.Conf) {
	this.key = section.String("key", "")
	this.action = section.String("act", "")
	this.cur = section.String("cur", "")
	this.rang = section.IntList("range", nil)
}

type EsFilter struct {
	sink       int
	converters []converter
}

func (this *EsFilter) Init(config *conf.Conf) {
	const CONV = "converts"
	this.sink = config.Int("sink", 0)
	for i := 0; i < len(config.List(CONV, nil)); i++ {
		section, err := config.Section(fmt.Sprintf("%s[%d]", CONV, i))
		if err != nil {
			panic(err)
		}

		c := converter{}
		c.load(section)
		this.converters = append(this.converters, c)
	}
}

func (this *EsFilter) Run(r engine.FilterRunner, e *engine.EngineConfig) error {
	globals := engine.Globals()
	if globals.Verbose {
		globals.Println("EsFilter started")

	}

	geodbFile := e.String("geodbfile", "")
	if err := als.LoadGeoDb(geodbFile); err != nil {
		panic(err)
	}
	if globals.Verbose {
		globals.Printf("Load geodb from %s\n", geodbFile)
	}

	var (
		pack   *engine.PipelinePack
		ok     = true
		inChan = r.InChan()
	)

	for ok && !globals.Stopping {
		select {
		case pack, ok = <-inChan:
			if !ok {
				break
			}

			if this.handlePack(pack) {
				r.Inject(pack)
			} else {
				pack.Recycle()
			}
		}
	}

	return nil
}

func (this *EsFilter) handlePack(pack *engine.PipelinePack) bool {
	pack.Message.SetField("area", pack.Message.Area)
	pack.Message.SetField("ts", pack.Message.Timestamp)
	pack.Message.Sink = this.sink

	for _, conv := range this.converters {
		switch conv.action {
		case "money":
			amount, err := pack.Message.FieldValue(conv.key, als.KEY_TYPE_MONEY)
			if err != nil {
				// has no such field
				continue
			}

			currency, err := pack.Message.FieldValue(conv.cur, als.KEY_TYPE_STRING)
			if err != nil {
				// has money field, but no currency field?
				return false
			}

			pack.Message.SetField("usd",
				als.MoneyInUsdCents(currency.(string), amount.(int)))

		case "ip":
			ip, err := pack.Message.FieldValue(conv.key, als.KEY_TYPE_IP)
			if err != nil {
				continue
			}

			pack.Message.SetField("cntry", als.IpToCountry(ip.(string)))

		case "range":
			if len(conv.rang) < 2 {
				continue
			}

			val, err := pack.Message.FieldValue(conv.key, als.KEY_TYPE_INT)
			if err != nil {
				continue
			}

			pack.Message.SetField(conv.key+"_rg", als.GroupInt(val.(int), conv.rang))

		case "del":
			pack.Message.DelField(conv.key)
		}
	}

	return true
}

func init() {
	engine.RegisterPlugin("EsFilter", func() engine.Plugin {
		return new(EsFilter)
	})
}