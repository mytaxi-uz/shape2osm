package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func BuildingAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	tags = append(tags, osm.Tag{
		Key:   "building",
		Value: "yes",
	})

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "STOREY":
			key = "building:levels"
			value = strings.Split(attr, ".")[0]
		case "HOUSE_NUM":
			key = "addr:housenumber"
			value = attr
		}

		if key != "" {
			tags = append(tags, osm.Tag{
				Key:   key,
				Value: value,
			})
		}
	}

	return
}
