package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func WaterAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value string

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			key = "water"
			switch attr {
			case "62":
				value = "reservoir"
				tags = append(tags, osm.Tag{
					Key:   "natural",
					Value: "water",
				})
			case "61", "68":
				value = "river"
				tags = append(tags, osm.Tag{
					Key:   "natural",
					Value: "water",
				})
			case "64", "74", "75":
				value = "canal"
				tags = append(tags, osm.Tag{
					Key:   "natural",
					Value: "water",
				})
			case "65", "271", "72", "73":
				key = "leisure"
				value = "swimming_pool"
			default:
				key = ""
			}
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
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
