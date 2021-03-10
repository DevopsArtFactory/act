package builder

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/act/pkg/tools"
)

type Flags struct {
	Region   string `json:"region"`
	Duration int    `json:"duration"`
	Profile  string `json:"profile"`
}

func ParseFlags() (*Flags, error) {
	keys := viper.AllKeys()
	flags := Flags{}

	val := reflect.ValueOf(&flags).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		key := strings.ReplaceAll(typeField.Tag.Get("json"), "_", "-")
		if tools.IsStringInArray(key, keys) {
			t := val.FieldByName(typeField.Name)
			if t.CanSet() {
				switch t.Kind() {
				case reflect.String:
					t.SetString(viper.GetString(key))
				case reflect.Int:
					t.SetInt(viper.GetInt64(key))
				case reflect.Bool:
					t.SetBool(viper.GetBool(key))
				}
			}
		}
	}

	return &flags, nil
}
