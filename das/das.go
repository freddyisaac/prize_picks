package das

import (
	"context"
	"fmt"
	"log"
	"strings"

	"database/sql"
	_ "github.com/lib/pq"
)

// poorly organised monolithic data access interface
// should make more modular
type DataAccessProvider interface {
	NewCage(ctx context.Context, cap int, kind string) (int, error)
	AddDinosaur(ctx context.Context, d Dinosaur) error
	PlaceDinosaurInCage(ctx context.Context, cageID int, d Dinosaur) error
	GetCages(ctx context.Context, optStatus ...string) ([]Cage, error)
	GetDinosaursForCage(ctx context.Context, cageID int) ([]Dinosaur, error)
	GetDinosaurs(ctx context.Context, opts ...string) ([]Dinosaur, error)
	SetCageStatus(ctx context.Context, cageID int, status string) error
	Close()
}

type PsqlDataProvider struct {
	db *sql.DB
}

// connect to database and return a data access object
func Connect(host, port, user, pwd, dbname string) (DataAccessProvider, error) {
	psqlConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbname)
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PsqlDataProvider{db: db}, nil
}

func (pdb *PsqlDataProvider) Close() {
	pdb.db.Close()
}

// place dinosaur in a given cage
// this does not enforce constraints correctly and my be subject to race condition
// leading to an overly cramped cage
// probably requires row level locking and transactional integrity
func (pdb *PsqlDataProvider) PlaceDinosaurInCage(ctx context.Context, cageID int, d Dinosaur) error {
	err := pdb.CheckCage(ctx, cageID, d.Diet)
	if err != nil {
		return err
	}
	sqlStmt := `INSERT INTO dinosaurs (species, name, diet, cage) VALUES ($1, $2, $3, $4)`
	_, err = pdb.db.Exec(sqlStmt, strings.ToLower(d.Species), strings.ToLower(d.Name), d.Diet, cageID)
	return err
}

func (pdb *PsqlDataProvider) CheckCage(ctx context.Context, cageID int, diet string) error {
	sqlStmt := `UPDATE cages SET  count=count+1 WHERE id = $1 AND kind = $2 AND count < capacity`
	res, err := pdb.db.ExecContext(ctx, sqlStmt, cageID, diet)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 1 {
		return nil
	}
	return fmt.Errorf("unable to place dino in cage %d", cageID)
}

// TODO: MAP DB LEVEL ERRORS TO APP LEVEL ERRORS

func (pdb *PsqlDataProvider) AddDinosaur(ctx context.Context, d Dinosaur) error {
	// this could be a code const but is probably ok for the demo
	id, err := pdb.GetFreeCage(ctx, d.Diet)
	if err != nil {
		return err
	}
	// this is a replica of above - consolidate
	sqlStmt := `INSERT INTO dinosaurs (species, name, diet, cage) VALUES ($1, $2, $3, $4)`
	_, err = pdb.db.Exec(sqlStmt, strings.ToLower(d.Species), strings.ToLower(d.Name), d.Diet, id)
	return err
}

const (
	StatusDown   = "DOWN"
	StatusActive = "ACTIVE"
	CageCapacity = 20
)

func ValidStatus(status string) bool {
	switch status {
	case StatusDown, StatusActive:
		return true
	default:
		return false
	}
	return false
}

func (pdb *PsqlDataProvider) NewCage(ctx context.Context, cap int, kind string) (int, error) {
	if cap < 1 {
		return -1, fmt.Errorf("cage capacity < 1 not permitted")
	}
	sqlStmt := `INSERT INTO cages (status, capacity, count, kind) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := pdb.db.QueryRow(sqlStmt, StatusActive, cap, 0, kind).Scan(&id)
	return id, err
}

// get a free cage of the required type
// if none is available then create a new one
func (pdb *PsqlDataProvider) GetFreeCage(ctx context.Context, diet string) (int, error) {
	sqlStmt := `SELECT id FROM cages WHERE ( count < capacity ) AND status = "ACTIVE" AND kind = $1 FOR UPDATE`
	rows, err := pdb.db.QueryContext(ctx, sqlStmt, diet)
	if err != nil {
		log.Printf("GetFreeCage() : %v", err)
		return 0, err
	}
	var id int
	count := 0
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
		count++
	}
	if count == 0 {
		id, err = pdb.NewCage(ctx, CageCapacity, diet)
		if err != nil {
			return 0, err
		}
	}
	err = pdb.PlaceInCage(ctx, id)
	return id, nil
}

func (pdb *PsqlDataProvider) GetCages(ctx context.Context, optStatus ...string) ([]Cage, error) {
	var cages []Cage
	var opts []interface{}
	sqlStmt := `SELECT id, status, capacity, count, kind FROM cages`
	if len(optStatus) != 0 {
		sqlStmt = sqlStmt + ` WHERE status = $1`
		opts = append(opts, &optStatus[0])
	}
	rows, err := pdb.db.QueryContext(ctx, sqlStmt, opts...)
	if err != nil {
		return cages, err
	}
	defer rows.Close()
	for rows.Next() {
		cage := Cage{}
		err := rows.Scan(&cage.ID, &cage.Status, &cage.Capacity, &cage.Count, &cage.Kind)
		if err != nil {
			return cages, err
		}
		cages = append(cages, cage)
	}
	return cages, nil
}

func (pdb *PsqlDataProvider) PlaceInCage(ctx context.Context, id int) error {
	sqlStmt := `UPDATE cages SET count = count + 1 WHERE id = $1`
	_, err := pdb.db.ExecContext(ctx, sqlStmt, id)
	return err
}

func (pdb *PsqlDataProvider) GetDinosaursForCage(ctx context.Context, cageID int) ([]Dinosaur, error) {
	var dinos []Dinosaur
	sqlStmt := `SELECT id, species, name, diet, cage FROM dinosaurs WHERE cage=$1`
	rows, err := pdb.db.QueryContext(ctx, sqlStmt, cageID)
	if err != nil {
		return dinos, err
	}
	defer rows.Close()
	for rows.Next() {
		dino := Dinosaur{}
		err := rows.Scan(&dino.ID, &dino.Name, &dino.Diet, &dino.Cage)
		if err != nil {
			return dinos, err
		}
		dinos = append(dinos, dino)
	}
	return dinos, nil
}

func (pdb *PsqlDataProvider) GetDinosaurs(ctx context.Context, species ...string) ([]Dinosaur, error) {
	var dinos []Dinosaur
	var opts []interface{}
	sqlStmt := `SELECT id, species, name, diet, cage FROM dinosaurs`
	if len(species) != 0 {
		sqlStmt = sqlStmt + ` WHERE species = $1`
		opts = append(opts, &species[0])
	}
	log.Printf("sql : %v", sqlStmt)
	rows, err := pdb.db.QueryContext(ctx, sqlStmt, opts...)
	if err != nil {
		log.Printf("err : %v", err)
		return dinos, err
	}
	defer rows.Close()
	for rows.Next() {
		dino := Dinosaur{}
		err := rows.Scan(&dino.ID, &dino.Species, &dino.Name, &dino.Diet, &dino.Cage)
		if err != nil {
			return dinos, err
		}
		dinos = append(dinos, dino)
	}
	return dinos, nil
}

func (pdb *PsqlDataProvider) CheckCageDiet(ctx context.Context, cageID int, diet string) (bool, error) {
	sqlStmt := `SELECT count(*) FROM cage WHERE id = $1 AND diet=$2`
	var count int
	err := pdb.db.QueryRowContext(ctx, sqlStmt, cageID, diet).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func (pdb *PsqlDataProvider) SetCageStatus(ctx context.Context, cageID int, status string) error {
	sqlStmt := `UPDATE cages SET status = $1 WHERE id = $2`
	if status == StatusDown {
		// ensure that if we power down a cage then it must not be empty
		sqlStmt = sqlStmt + ` AND count = 0`
	}
	log.Printf("sql : %s", sqlStmt)
	resp, err := pdb.db.ExecContext(ctx, sqlStmt, status, cageID)
	if resp != nil {
		num, err := resp.RowsAffected()
		if err != nil {
			return err
		}
		if num == 0 {
			return fmt.Errorf("cage status unchanged")
		}
	}
	return err
}
