package oracle

import (
	"database/sql"
	"time"
	"flag"
	"os"
	"strconv"
	"log"
	"fmt"

	_ "github.com/godror/godror"

	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func connect (connectionString string) (*sql.DB) {

	db, err := sql.Open("godror", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func rowsCount (db *sql.DB, tablename string) (float64) {

	var rows_count float64
	query := "select count(*) from " + tablename
	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&rows_count)
	}

	return rows_count
}

func getenvStr(key string) (string) {

	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("getenv: Environment variable '%s' is empty", key)
	}
	return v
}

func getenvInt(key string) (int) {

	s := getenvStr(key)

	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func main() {

	listenPort := getenvStr("ListenPort")
	selectTimeout := getenvInt("SelectTimeout")
	oracleHost := getenvStr("OracleHost")
	oraclePort := getenvStr("OraclePort")
	oracleSystemName := getenvStr("OracleSystemName")
	oracleUser := getenvStr("OracleUser")
	oraclePassword := getenvStr("OraclePassword")
	oracleTableName := getenvStr("OracleTableName")

	connectionString := fmt.Sprintf("user=%s password=%s connectString=%s:%s/%s", oracleUser, oraclePassword, oracleHost, oraclePort, oracleSystemName)  

	addr := flag.String("listen-address", "0.0.0.0:" + listenPort,
 	 "The address to listen on for HTTP requests.")

	count := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "oracle",
			Name:      "table_rows",
			Help:      "Number of rows in oracle table",
		})
	prometheus.MustRegister(count)

	db := connect(connectionString)
	defer db.Close()
	
	go func() {
		for {
			count.Set(rowsCount(db, oracleTableName))
			time.Sleep(time.Duration(selectTimeout) * time.Second)
		}
	  }()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting web server at %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	  }
}