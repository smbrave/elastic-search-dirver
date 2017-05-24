# elastic-search-dirver
elastic-search的sql驱动，可通过sql.DB访问

```golang


package main

import (
	"database/sql"
	"fmt"

	_ "github.com/smbrave/elastic-search-driver/es-driver"
)

func main() {
	dsn := `{"username":"111","password":"123","addrlist":["10.0.54.127:9200","10.0.54.127:9200"]}`
	db, err := sql.Open("es", dsn)
	if err != nil {
		fmt.Printf("open:%s\n", err.Error())
		return
	}

	rows, err := db.Query("SELECT a,b FROM automarket_oil.order_info WHERE order_type=2 AND order_time>10 ORDER BY a desc,b ASC")
	if err != nil {
		fmt.Printf("query11:%s\n", err.Error())
		return
	}
	cols, _ := rows.Columns()
	fmt.Printf("%+v\n", cols)
	for rows.Next() {
		var aa1 string
		var bb string
		var cc string
		rows.Scan(&aa1, &bb, &cc)
		fmt.Printf("aa:%s, bb:%s,cc:%s\n", aa1, bb, cc)
	}
}


```