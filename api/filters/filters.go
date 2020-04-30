package filters

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gdey/forge-api-go-client/api"
)

const (
	keyFmt            = "filter[%v]%s"
	KeyExtensionType  = "filter[extension.type]"
	KeyHidden         = "filter[hidden]"
	KeyID             = "filter[id]"
	KeyMimeType       = "filter[mimeType]"
	KeyName           = "filter[name]"
	KeyType           = "filter[type]"
	KeyVersionNumber  = "filter[versionNumber]"
	KeyPageNumber     = "page[number]"
	KeyPageLimit      = "page[limit]"
	KeyPageStartAt    = "page[startAt]"
	KeyPageBeginsWith = "page[beginsWith]"
)

const (
	ComparisonNone           = Comparison(0)
	ComparisonEqual          = Comparison(1)
	ComparisonLess           = Comparison(2)
	ComparisonLessOrEqual    = Comparison(3)
	ComparisonGreater        = Comparison(4)
	ComparisonGreaterOrEqual = Comparison(5)
	ComparisonStartsWith     = Comparison(6)
	ComparisonEndsWith       = Comparison(7)
	ComparisonContains       = Comparison(8)
)

type Comparison uint

func (c Comparison) String() string {
	switch c {
	default:
		return ""
	case ComparisonEqual:
		return "-eq"
	case ComparisonLess:
		return "-lt"
	case ComparisonLessOrEqual:
		return "-le"
	case ComparisonGreater:
		return "-gt"
	case ComparisonGreaterOrEqual:
		return "-ge"
	case ComparisonStartsWith:
		return "-starts"
	case ComparisonEndsWith:
		return "-ends"
	case ComparisonContains:
		return "-contains"
	}
}

func RunAll(values url.Values, filters ...api.Filterer) (err error) {
	for _, flt := range filters {
		err = flt.Add(values)
		if err != nil {
			return err
		}
	}
	return nil
}

func addStringSlice(key string, strs []string, values url.Values) error {
	for _, str := range strs {
		values.Add(key, str)
	}
	return nil
}

type Page struct {
	Limit      uint8
	Number     uint
	BeginsWith string
	StartAt    string
}

func (filter *Page) Add(values url.Values) error {
	if filter == nil {
		return nil
	}
	if filter.BeginsWith != "" {
		values.Add(KeyPageBeginsWith, filter.BeginsWith)
	}
	if filter.StartAt != "" {
		values.Add(KeyPageStartAt, filter.StartAt)
	}
	if filter.Number != 0 {
		values.Add(KeyPageNumber, strconv.FormatUint(uint64(filter.Number), 10))
	}
	limit := filter.Limit
	if filter.Limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	values.Add(KeyPageLimit, strconv.FormatUint(uint64(limit), 10))
	return nil
}

type QueryParam struct {
	Key   string
	Value string
}

func (query QueryParam) Add(values url.Values) error {
	values.Add(query.Key, query.Value)
	return nil
}

type Filter struct {
	Name       string
	Values     []string
	Comparison Comparison
}

func (filter Filter) Add(values url.Values) error {
	key := fmt.Sprintf(keyFmt, filter.Name, filter.Comparison)
	return addStringSlice(key, filter.Values, values)
}

type VersionNumbers []VersionNumber

func (filter VersionNumbers) Add(values url.Values) (err error) {
	if len(filter) == 0 {
		return nil
	}
	for i := range filter {
		if err = filter[i].Add(values); err != nil {
			return err
		}
	}
	return nil
}

type VersionNumber struct {
	Comparison Comparison
	Version    int64
}

func (filter VersionNumber) Add(values url.Values) error {
	key := fmt.Sprintf(keyFmt, "versionNumber", filter.Comparison)
	values.Add(key, strconv.FormatInt(filter.Version, 10))
	return nil
}

type Hidden uint8

const (
	ShowHidden    = Hidden(1)
	ShowNonHidden = Hidden(2)
)

func (filter Hidden) Add(values url.Values) error {
	if filter&ShowHidden == ShowHidden {
		values.Add(KeyHidden, "true")
	}
	if filter&ShowNonHidden == ShowNonHidden {
		values.Add(KeyHidden, "false")
	}
	return nil
}

//Type is used to filter by the type of the ref target
type Type []string

func (filter Type) Add(values url.Values) error {
	return addStringSlice(KeyType, []string(filter), values)
}

// ID is used to filter by the id of the ref target
type ID []string

func (filter ID) Add(values url.Values) error {
	return addStringSlice(KeyID, []string(filter), values)
}

// Name is used to filter by the name of the ref target
type Name []string

func (filter Name) Add(values url.Values) error {
	return addStringSlice(KeyName, []string(filter), values)
}

// ExtensionType is used to filter by the extension type
type ExtensionType []string

func (filter ExtensionType) Add(values url.Values) error {
	return addStringSlice(KeyExtensionType, []string(filter), values)
}
