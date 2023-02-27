package snowflake

import (
	"bluebell/settings"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(nfg *settings.NameConfig) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", nfg.StartTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(nfg.MachineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
