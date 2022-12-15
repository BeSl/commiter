package database

func SchemaUser() string {
	return `CREATE TABLE users (
    		id serial PRIMARY KEY,
			extId UUID, 
			name text,
    		is_admin boolean,
    		gitlogin text,
			tgID text
		);`
}

func SchemaEProc() string {
	return `CREATE TABLE extprocVersion (
    		id serial PRIMARY KEY,
			authorversion UUID,
			extId UUID, 
			name text,
    		BinaryData text,
    		Filename text,
			Processed boolean
		);`
}
