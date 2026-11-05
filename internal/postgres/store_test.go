package postgres

import (
	"context"
	"os"
	"testing"
)

func TestMigrationsAreRepeatableAndSchemaIsReady(t *testing.T) {
	url := os.Getenv("SMON_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set SMON_TEST_DATABASE_URL to a dedicated empty test database")
	}
	ctx := context.Background()
	s, err := Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for i := 0; i < 2; i++ {
		if err := s.Migrate(ctx); err != nil {
			t.Fatalf("migration run %d: %v", i, err)
		}
	}
	if err := s.Ready(ctx); err != nil {
		t.Fatal(err)
	}
	// Pending must be representable without manufacturing a successful result.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	var host, command, service int64
	if err := tx.QueryRow(ctx, "INSERT INTO hosts(name,address) VALUES ('smon-test-host','127.0.0.1') RETURNING id").Scan(&host); err != nil {
		t.Fatal(err)
	}
	if err := tx.QueryRow(ctx, "INSERT INTO commands(name,executable) VALUES ('smon-test-check','/bin/true') RETURNING id").Scan(&command); err != nil {
		t.Fatal(err)
	}
	if err := tx.QueryRow(ctx, "INSERT INTO services(host_id,command_id,description) VALUES ($1,$2,'test') RETURNING id", host, command).Scan(&service); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := tx.QueryRow(ctx, "SELECT count(*) FROM service_states WHERE service_id=$1", service).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("new service must remain pending")
	}
	_, err = tx.Exec(ctx, "INSERT INTO service_states(service_id,state,state_type,attempt,output,checked_at,state_changed_at) VALUES ($1,4,'HARD',1,'invalid',now(),now())", service)
	if err == nil {
		t.Fatal("invalid plugin state was accepted")
	}
}
