package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

const CassandraHost = "10.48.3.21"
const CassandraAuthUser = "binu.b.varghese"
const CassandraAuthPass = "Accenture01!"
const CassandraKeyspace = "test"
const CassandraPort = 9042

/*
use test;
CREATE TABLE patients (
id UUID,
firstname text,
lastname text,
age int,
email text,
city text,
PRIMARY KEY (id)
);
*/

// Session holds our connection to Cassandra
var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster(CassandraHost)
	pass := gocql.PasswordAuthenticator{CassandraAuthUser, CassandraAuthPass}
	cluster.Keyspace = CassandraKeyspace
	cluster.Authenticator = pass
	cluster.Consistency = gocql.One
	cluster.Port = CassandraPort
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}
