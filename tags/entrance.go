package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func EntranceAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	tags = append(tags, osm.Tag{
		Key:   "entrance",
		Value: "yes",
	})

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		field := strings.ToUpper(f.String())

		switch field {
		case "ENTRANCE":
			if attr != "0" && attr != "" {
				tags = append(tags, osm.Tag{
					Key:   "name",
					Value: attr,
				})
			}
		}
	}

	return
}
