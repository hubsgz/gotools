package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
	"io/ioutil"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

/**
 *  sql数据导出工具
 *	功能: 可通过配置sql语句导出大批量的数据,避免直接执行sql导出可能导致数据库负载高的问题
 *  用法参数:
      ./exporttool -d directory [-s size]
 *    -d  配置目录, 非必传, 默认为当前目录, 该目录必须有以下三个文件
			dsn.conf  配置数据库连接
			count.sql 统计总数据量sql
			list.sql  导出数据列表sql
       -s  分页size, 每页查询多少条数据， 非必传, 默认为10,  该值越高， 对数据库压力越大
     最终数据会导出到 -d 配置的目录下， 文件名为out-当前时间戳名.csv命名
 */

var cfgdir string
var size int64

func main()  {
	flag.StringVar(&cfgdir, "d", ".", "config dir")  //默认当前目录
	flag.Int64Var(&size, "s", 10, "page size")       //默认每次查询10条记录
	flag.Parse()
	log.Println("use cfgdir ", cfgdir)

	sqlfile := fmt.Sprintf("%s%slist.sql", cfgdir, string(os.PathSeparator))
	sqlfile_count := fmt.Sprintf("%s%scount.sql", cfgdir, string(os.PathSeparator))
	dbfile := fmt.Sprintf("%s%sdsn.conf", cfgdir, string(os.PathSeparator))
	tofile := fmt.Sprintf("%s%sout-%d.csv", cfgdir, string(os.PathSeparator), time.Now().Unix())
	log.Println(tofile)

	db := initdb(dbfile)
	defer db.Close()

	sqlcount, err := ioutil.ReadFile(sqlfile_count)
	if err != nil {
		log.Panic(err)
	}
	sqlcount_s := string(sqlcount)
	fmt.Println(sqlcount_s)

	sqlf, err := ioutil.ReadFile(sqlfile)
	if err != nil {
		log.Panic(err)
	}
	sql_s := string(sqlf)
	fmt.Println(sql_s)

	outfile, err := os.Create(tofile)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	outfile.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	// csv文件writer
	w := csv.NewWriter(outfile)

	var num int64
	db.QueryRow(sqlcount_s).Scan(&num)
	log.Println("total count=", num)

	pagecount := int64(math.Ceil(float64(num) / float64(size)))
	log.Println("pagecount=", pagecount)

	var i int64
	wheader := false
	for i=1; i<=pagecount; i++  {
		offset := (i-1) * size
		pagesql := fmt.Sprintf("%s limit %d,%d", sql_s, offset, size)
		//log.Println(pagesql)
		rows, err := db.Query(pagesql)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}
		columns,err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}

		row := make([]sql.RawBytes, len(columns))
		row_p := make([]interface{}, len(columns))
		for ik := range row_p {
			row_p[ik] = &row[ik]
		}
		if !wheader {
			w.Write(columns)
			wheader = true
		}
		for rows.Next() {
			err = rows.Scan(row_p...)
			if err != nil {
				log.Fatal("scan err=",err)
			}
			//log.Println(row)
			var val string
			var arr = make([]string, len(columns))
			for k,v := range row {
				if v == nil {
					val = ""
				} else {
					val = string(v)
				}
				//log.Println(columns[k],"=",val)
				arr[k] = val
			}
			w.Write(arr)
			//break;
		}

		log.Printf("page=%d/%d", i, pagecount)
		if i > 10 {
			//break;
		}
		if i%10 == 0 {
			w.Flush()
		}
	}

	w.Flush()
	log.Println("all finish")

	var in string
	fmt.Printf("按回车退出: ")
	fmt.Scanln(&in)
	fmt.Printf("%s\n", in)
}

func initdb(dbfile string) *sql.DB {
	var err error
	dsn, err := ioutil.ReadFile(dbfile)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(dsn))
	db, err := sql.Open("mysql", string(dsn))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	return db
}
