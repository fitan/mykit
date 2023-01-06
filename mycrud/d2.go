package mycrud

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm/schema"
	"net/http"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"sort"
	"strings"
)

func (c *Core) D2Handler(m *mux.Router) {
	m.HandleFunc("/d2/tables", func(writer http.ResponseWriter, request *http.Request) {
		var ss []schema.Schema
		for _, v := range c.tables {
			ss = append(ss, v.schema)
		}

		s := GenD2(ss)
		writer.Write(D2(s))
	})

	m.HandleFunc("/d2/tables/{tableName}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		tableName := vars["tableName"]
		msg, err := c.tableMsg(tableName)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		var ss []schema.Schema
		for _, v := range msg.schema.Relationships.Relations {
			ss = append(ss, *(v.FieldSchema))
		}
		ss = append(ss, msg.schema)
		s := GenD2(ss)
		writer.Write(D2(s))
	})
}

func D2(s string) []byte {
	ruler, _ := textmeasure.NewRuler()
	defaultLayout := func(ctx context.Context, g *d2graph.Graph) error {
		return d2dagrelayout.Layout(ctx, g, nil)
	}
	diagram, _, _ := d2lib.Compile(context.Background(), s, &d2lib.CompileOptions{
		Layout:  defaultLayout,
		Ruler:   ruler,
		ThemeID: d2themescatalog.GrapeSoda.ID,
	})
	out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad:    d2svg.DEFAULT_PADDING,
		Sketch: false,
	})
	return out
	//_ = ioutil.WriteFile(filepath.Join("out.svg"), out, 0600)
}

func GenD2(ss []schema.Schema) string {
	var d2 []string
	var relation []string
	relationRecord := make(map[string]struct{})

	for _, v := range ss {
		for _, value := range v.Relationships.Relations {
			relationV := v.Table + " -> " + value.FieldSchema.Table + ": " + string(value.Type)
			if _, ok := relationRecord[relationV]; ok {
				continue
			}

			if v.Table != value.FieldSchema.Table {

				relationRecord[relationV] = struct{}{}
				relation = append(relation, relationV)
			}
		}

		var fs sqlTableFields
		for _, field := range v.FieldsByName {
			t := field.FieldType.Kind().String()
			_, pk := field.TagSettings["PRIMARYKEY"]
			comment := field.TagSettings["COMMENT"]
			fs = append(fs, sqlTableField{
				Name: field.Name,
				T:    fmt.Sprintf("%s (%s)", t, comment),
				Pk:   pk,
			})
		}

		sort.Sort(fs)

		d2 = append(d2, sqlTable(v.Table, fs))
	}

	sort.Sort(sort.StringSlice(d2))

	d2 = append(d2, relation...)

	return strings.Join(d2, "\n")
}

type sqlTableField struct {
	Name string
	T    string
	Pk   bool
}

type sqlTableFields []sqlTableField

func (s sqlTableFields) Len() int {
	return len(s)
}

func (s sqlTableFields) Less(i, j int) bool {
	if s[i].Name < s[j].Name {
		return true
	}
	return false
}

func (s sqlTableFields) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sqlTable(tableName string, fs []sqlTableField) string {
	f := `
%s: {
  shape: sql_table
  %s
}`
	var fields []string
	for _, v := range fs {
		var pkS string
		if v.Pk {
			pkS = `{constraint: primary_key}`
		}
		fields = append(fields, v.Name+": "+v.T+" "+pkS)
	}

	return fmt.Sprintf(f, tableName, strings.Join(fields, "  \n"))
}
