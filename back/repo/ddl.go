package repo

import "fmt"

var ddls = []func(){
	createTaskTable,
}

func createTaskTable() {
	conn := getDbMust()
	defer conn.Close()

	if _, err := conn.Exec(`
	create table if not exists task (
		id integer primary key not null,
		owner_name text not null,
		title text not null,
		live_id text unique not null,
		live_time timestamp,
		status text not null,
		error_info text not null,
		details text not null,

		created timestamp not null default current_timestamp,
		updated timestamp
	);`); err != nil {
		panic(err)
	}

	addUpdatedTriggerOnTable("task")
}

func addUpdatedTriggerOnTable(table string) {
	conn := getDbMust()
	defer conn.Close()
	setTriggerSql := fmt.Sprintf(`
	CREATE TRIGGER IF NOT EXISTS update_timestamp
	AFTER UPDATE ON %s
	FOR EACH ROW
	BEGIN
		UPDATE %s SET updated = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;`, table, table)
	if _, err := conn.Exec(setTriggerSql); err != nil {
		panic(err)
	}
}
