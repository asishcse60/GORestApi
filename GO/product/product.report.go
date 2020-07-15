package product

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type ReportFilter struct {
	NameFilter string `json:"productName"`
	ManufactureFilter string `json:"manufacturer"`
	SKUFilter string `json:"sku"`
}

func handleProductReport(w http.ResponseWriter, r *http.Request)  {
	switch r.Method {
		case http.MethodPost:
			var productReport ReportFilter
			err:=json.NewDecoder(r.Body).Decode(&productReport)
			if err != nil{
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			products,err:=searchForProductData(productReport)
			if err != nil{
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tem:=template.New("report.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool {return i%x==0}})
			tem,err = tem.ParseFiles(path.Join("templates", "report.gotmpl"))
			if err != nil{
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var tpl bytes.Buffer
		//	var product Product
			if len(products) > 0{
				//product = products[0]
				err=tem.Execute(&tpl, products)
			}else{
				w.WriteHeader(http.StatusNotFound)
				return
			}
			rdr:=bytes.NewReader(tpl.Bytes())
			w.Header().Set("Content-Disposition", "Attachment")
			http.ServeContent(w, r, "report.html",time.Now(), rdr)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
