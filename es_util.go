package es_driver

import (
	"errors"
	"fmt"
	"strings"
)

func parseESSQL(query string) (*ESConfSQL, error) {
	// SELECT a,bc as b,c from tb where a=b and b=a order by a desc
	// INSERT INTO Store_Information (Store_Name,Sales,Txn_Date) VALUES ('Los Angeles',900,'Jan-10-1999');
	// DELETE FROM 表名称 WHERE 列名称 = 值

	query = strings.Trim(query, "\t\n ")
	sql := new(ESConfSQL)
	pos := strings.Index(query, " ")
	if pos == -1 {
		return nil, errors.New("sql error")
	}

	sql.Method = strings.ToUpper(query[:pos])
	var err error
	if sql.Method == "SELECT" {
		err = parseSELECT(query, sql)
	} else if sql.Method == "INSERT" {
		err = parseINSERT(query, sql)
	} else if sql.Method == "DELETE" {
		err = parseDELETE(query, sql)
	} else {
		return nil, fmt.Errorf("no support [%s]", sql.Method)
	}

	return sql, err
}

// SELECT a,bc as b,c from tb where a=b and b=a order by a desc
func parseSELECT(query string, sql *ESConfSQL) error {
	items := strings.Split(query, " ")
	fields := make([]string, 0)
	where := make([]string, 0)
	order := make([]string, 0)
	var isField = false
	var isTable = false
	var isWhere = false
	var isOrder = false

	idx := 1
	for {
		if idx >= len(items) {
			break
		}

		if !isField {
			if strings.ToUpper(items[idx]) == "FROM" {
				isField = true
				idx += 1
				continue
			}
			fields = append(fields, items[idx])
			idx += 1
			continue
		}

		if !isTable {
			fs := strings.Split(items[idx], ".")
			if len(fs) != 2 {
				return errors.New("must form database.table")
			}
			sql.Database = fs[0]
			sql.Table = fs[1]
			isTable = true
			idx += 1
			continue
		}

		if !isWhere {
			if strings.ToUpper(items[idx]) == "WHERE" {
				isWhere = true
				for {
					idx += 1
					if idx >= len(items) || strings.ToUpper(items[idx]) == "ORDER" {
						break
					}
					if strings.ToUpper(items[idx]) != "AND" {
						where = append(where, items[idx])
					}
				}
				continue
			}
		}

		if !isOrder {
			if strings.ToUpper(items[idx]) == "ORDER" {
				idx += 1
				if idx >= len(items) {
					break
				}
				if strings.ToUpper(items[idx]) == "BY" {
					isOrder = true
					for {
						idx += 1
						if idx >= len(items) {
							break
						}
						order = append(order, items[idx])
					}
				}
			}
		}

	}

	fields = strings.Split(strings.Join(fields, " "), ",")
	for i := 0; i < len(fields); i++ {
		f := strings.Trim(fields[i], "\t\n ")
		fs := strings.Split(f, " ")
		if len(fs) == 1 {
			sql.Field = append(sql.Field, f)
			sql.Alias = append(sql.Alias, f)
		} else if len(fs) == 3 && strings.ToUpper(fs[1]) == "AS" {
			sql.Field = append(sql.Field, fs[0])
			sql.Alias = append(sql.Alias, fs[2])
		} else {
			return fmt.Errorf("%s error", f)
		}
	}

	order = strings.Split(strings.Join(order, " "), ",")
	for i := 0; i < len(order); i++ {
		sql.Order = append(sql.Order, strings.Trim(order[i], "\t\n "))
	}

	sql.Where = where
	return nil
}

func parseDELETE(query string, sql *ESConfSQL) error {
	return nil
}

func parseINSERT(query string, sql *ESConfSQL) error {
	return nil
}

// {"from":0,"post_filter":{"bool":{"must":[{"term":{"order_type":"2"}}]}},"size":10,"sort":[{"order_time":{"order":"desc"}}]}
func parseFilterQuery(filters []string) (interface{}, error) {

	boolQuery := NewBoolQuery()

	for _, filter := range filters {
		filter := strings.Trim(filter, "\r\n\t ")
		var pos int

		pos = strings.Index(filter, ">=")
		if pos != -1 {
			rangeQuery := NewRangeQuery(filter[:pos])
			rangeQuery.Gte(filter[pos+2:])
			boolQuery.Must(rangeQuery)
			continue
		}

		pos = strings.Index(filter, "<=")
		if pos != -1 {
			rangeQuery := NewRangeQuery(filter[:pos])
			rangeQuery.Lte(filter[pos+2:])
			boolQuery.Must(rangeQuery)
			continue
		}

		pos = strings.Index(filter, "=")
		if pos != -1 {
			rangeQuery := NewTermQuery(filter[:pos], filter[pos+1:])
			boolQuery.Must(rangeQuery)
			continue
		}

		pos = strings.Index(filter, ">")
		if pos != -1 {
			rangeQuery := NewRangeQuery(filter[:pos])
			rangeQuery.Gt(filter[pos+1:])
			boolQuery.Must(rangeQuery)
			continue
		}

		pos = strings.Index(filter, "<")
		if pos != -1 {
			rangeQuery := NewRangeQuery(filter[:pos])
			rangeQuery.Lt(filter[pos+1:])
			boolQuery.Must(rangeQuery)
			continue
		}

	}
	return boolQuery.Source()

}
