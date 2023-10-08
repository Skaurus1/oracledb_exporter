# oracledb_exporter

Counts the number of rows in the table and sends to the Prometheus endpoint

Build:
go get oracle
go build -o oracle

Requires GCC to build and Oracle Instant Client, libaio1 on the system to run
https://www.oracle.com/cis/database/technologies/instant-client/downloads.html


Requires environment variables:

ListenPort for prometheus endpoint
SelectTimeout from Oracle DB
OracleHost
OraclePort
OracleSystemName
OracleUser
OraclePassword
OracleTableName
