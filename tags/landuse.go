package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func LanduseAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			switch attr {
			case "13":
				tags = append(tags, osm.Tag{
					Key:   "boundary",
					Value: "administrative",
				})
				tags = append(tags, osm.Tag{
					Key:   "admin_level",
					Value: "2",
				})
			case "9":
				tags = append(tags, osm.Tag{
					Key:   "boundary",
					Value: "administrative",
				})
				tags = append(tags, osm.Tag{
					Key:   "admin_level",
					Value: "4",
				})
			case "103", "111":
				key = "leisure"
				value = "park"
			case "92":
				key = "landuse"
				value = "grass"
			case "340", "349", "350":
				key = "amenity"
				value = "hospital"
			case "207", "213", "216", "201":
				key = "landuse"
				value = "university"
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
			tag := osm.Tag{
				Key:   key,
				Value: value,
			}
			tags = append(tags, tag)
		}
	}

	return
}
