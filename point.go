package main

import (
	"reflect"
	"strings"

	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// poi, place

func convertPointToOSMNode(shapeReader *shp.Reader, shapeType string) {
	// fields from the attribute table (DBF)
	fields := shapeReader.Fields()

	t := reflect.TypeOf(&shp.Point{}).Elem()

	// loop through all features in the shapefile
	for shapeReader.Next() {
		num, p := shapeReader.Shape()
		if reflect.TypeOf(p).Elem() != t {
			continue
		}
		point := p.(*shp.Point)
		osmID++
		var tags osm.Tags

		switch shapeType {
		case "poi":
			tags = convertPoiAttrToOSMTag(num, fields, shapeReader)
		case "place":
			tags = convertPointPlaceAttrToOSMTag(num, fields, shapeReader)
		}
		node := osm.Node{
			ID:        osmID,
			Lat:       util.TruncateFloat64(point.Y),
			Lon:       util.TruncateFloat64(point.X),
			Tags:      tags,
			Version:   1,
			Timestamp: nowTime,
		}
		osmOut.Nodes = append(osmOut.Nodes, &node)
		if osmOut.Bounds.MinLat > point.Y {
			osmOut.Bounds.MinLat = point.Y
		}
		if osmOut.Bounds.MaxLat < point.Y {
			osmOut.Bounds.MaxLat = point.Y
		}
		if osmOut.Bounds.MinLon > point.X {
			osmOut.Bounds.MinLon = point.X
		}
		if osmOut.Bounds.MaxLon < point.X {
			osmOut.Bounds.MaxLon = point.X
		}
	}
}

func convertPoiAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value, str_type_uz, str_type_ru, str_type_en string

	for k, f := range fields {
		field := strings.ToUpper(f.String())
		switch field {
		case "TYP_STR_UZ":
			str_type_uz = reader.ReadAttribute(num, k)
		case "TYP_STR_RU":
			str_type_ru = reader.ReadAttribute(num, k)
		case "TYP_STR_EN":
			str_type_en = reader.ReadAttribute(num, k)
		}
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		case "STREET_UZ":
			key = "addr:street:uz"
			value = attr
			if str_type_uz != "" {
				value += " " + str_type_uz
			}
		case "STREET_RU":
			key = "addr:street:ru"
			value = attr
			if str_type_ru != "" {
				value += " " + str_type_ru
			}
		case "STREET_EN":
			key = "addr:street:en"
			value = attr
			if str_type_en != "" {
				value += " " + str_type_en
			}
		case "ADDRESS_UZ":
			key = "addr:housenumber"
			value = attr
		// case "CLASS":
		// key = "amenity"
		// value = strings.ToLower(attr)
		case "ADDRESS":
			key = "addr:full"
			value = attr
		case "PHONE":
			key = "phone"
			value = attr
		case "WORKTIME":
			key = "opening_hours"
			value = attr
		case "HTTP":
			key = "website"
			value = attr
		case "EMAIL":
			key = "email"
			value = attr
		case "ZIPCODE":
			key = "addr:postcode"
			value = attr

		case "TYP_COD":
			switch attr {
			case "371":
				key = "aeroway"
				value = "aerodrome"
			case "727", "245", "248":
				key = "shop"
				value = "convenience"
			case "728":
				key = "shop"
				value = "clothes"
			case "734":
				key = "shop"
				value = "furniture"
			case "729":
				key = "shop"
				value = "hardware"
			case "171":
				key = "amenity"
				value = "bank"
			case "630":
				key = "amenity"
				value = "bar"
			case "613":
				key = "amenity"
				value = "hairdresser"
			case "243":
				key = "amenity"
				value = "cafe"
			case "609":
				key = "shop"
				value = "car"
			case "366":
				key = "shop"
				value = "car_repair"
			case "301":
				key = "amenity"
				value = "place_of_worship"
				tag := osm.Tag{
					Key:   "religion",
					Value: "muslim",
				}
				tags = append(tags, tag)
			case "303", "304", "622":
				key = "amenity"
				value = "place_of_worship"
				tag := osm.Tag{
					Key:   "religion",
					Value: "christian",
				}
				tags = append(tags, tag)
			case "327":
				key = "amenity"
				value = "cinema"
			case "213":
				key = "amenity"
				value = "college"
			case "615":
				key = "amenity"
				value = "dentist"
			case "362", "610":
				key = "amenity"
				value = "fuel"
			case "616", "624", "621":
				key = "office"
				value = "government"
			case "340", "350", "349":
				key = "amenity"
				value = "hospital"
			case "204":
				key = "amenity"
				value = "kindergarten"
			case "150":
				key = "tourism"
				value = "hotel"
			case "601", "250":
				key = "shop"
				value = "mall"
			case "333":
				key = "historic"
				value = "monument"
			case "329":
				key = "tourism"
				value = "museum"
			case "619", "709":
				key = "office"
				value = "company"
			case "103":
				key = "leisure"
				value = "park"
			case "240":
				key = "amenity"
				value = "pharmacy"
			case "247":
				key = "amenity"
				value = "restaurant"
			case "217", "218":
				key = "amenity"
				value = "school"
			case "733":
				key = "shop"
				value = "alcohol"
			case "731":
				key = "shop"
				value = "baby_goods"
			case "201":
				key = "amenity"
				value = "university"
			case "420", "440":
				key = "railway"
				value = "subway_entrance"
			case "325":
				key = "tourism"
				value = "zoo"
			case "271":
				key = "amenity"
				value = "swimming_pool"
			}
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

func convertPointPlaceAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value, str_type_uz, str_type_ru, str_type_en string

	for k, f := range fields {
		field := strings.ToUpper(f.String())
		switch field {
		case "TYP_STR_UZ":
			str_type_uz = reader.ReadAttribute(num, k)
		case "TYP_STR_RU":
			str_type_ru = reader.ReadAttribute(num, k)
		case "TYP_STR_EN":
			str_type_en = reader.ReadAttribute(num, k)
		}
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			key = "place"
			switch attr {
			case "20":
				tag := osm.Tag{
					Key:   "admin_level",
					Value: "6",
				}
				tags = append(tags, tag)
				value = "city"
			case "10":
				value = "town"
			case "11":
				value = "village"
			case "134", "130", "723":
				value = "neighbourhood"
			case "12":
				value = "district"
			case "21":
				value = "city"
			case "9":
				value = "region"
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
		case "STREET_UZ":
			key = "addr:street:uz"
			value = attr
			if str_type_uz != "" {
				value += " " + str_type_uz
			}
		case "STREET_RU":
			key = "addr:street:ru"
			value = attr
			if str_type_ru != "" {
				value += " " + str_type_ru
			}
		case "STREET_EN":
			key = "addr:street:en"
			value = attr
			if str_type_en != "" {
				value += " " + str_type_en
			}
		case "ADDRESS_UZ":
			key = "addr:housenumber"
			value = attr
		case "CLASS":
			key = "amenity"
			value = strings.ToLower(attr)
		case "ADDRESS":
			key = "addr:full"
			value = attr
		case "PHONE":
			key = "phone"
			value = attr
		case "WORKTIME":
			key = "opening_hours"
			value = attr
		case "HTTP":
			key = "website"
			value = attr
		case "EMAIL":
			key = "email"
			value = attr
		case "ZIPCODE":
			key = "addr:postcode"
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
