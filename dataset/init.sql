CREATE TABLE dinosaurs (
	id serial,
	species TEXT NOT NULL,
	name TEXT,
	cage INTEGER,
	diet CHAR(1) NOT NULL
);	

CREATE INDEX idx_species ON dinosaurs(species);
CREATE INDEX idx_cage_id ON dinosaurs(cage);

CREATE TABLE cages (
	id serial,
	status TEXT NOT NULL,
	capacity integer NOT NULL,
	count integer NOT NULL,
	kind char(1) NOT NULL
);

