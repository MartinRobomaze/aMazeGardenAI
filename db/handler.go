package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseHandler struct {
	DriverName string
	User       string
	Password   string
	Database   string
}

var db *sql.DB
var err error

var stmtIns *sql.Stmt
var stmtOut *sql.Stmt
var stmtDel *sql.Stmt

func (dbReader DatabaseHandler) Begin() (err error) {
	dbOpenString := dbReader.User + ":" + dbReader.Password + "@/" + dbReader.Database
	db, err = sql.Open(dbReader.DriverName, dbOpenString)
	return err
}

func (dbReader DatabaseHandler) Write(plantName string, wateredSoilMoisture int) (err error) {
	if stmtIns, err = db.Prepare("INSERT INTO plants_data (name, watered_soil_moisture) VALUES (?, ?)"); err != nil {
		return err
	}

	if _ ,err = stmtIns.Exec(plantName, wateredSoilMoisture); err != nil {
		return err
	}

	if err = stmtIns.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) Update(plantId int, plantName string, wateredSoilMoisture int) (err error) {
	if stmtIns, err = db.Prepare("UPDATE plants_data SET plant_name = ?, watered_soil_moisture = ? WHERE plant_id = ?"); err != nil {
		return err
	}

	if _ ,err = stmtIns.Exec(plantName, wateredSoilMoisture, plantId); err != nil {
		return err
	}

	if err = stmtIns.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) GetWateredSoilMoistureFromName(plantName string) (wateredSoilMoisture int, err error) {
	if stmtOut, err = db.Prepare("SELECT watered_soil_moisture FROM plants_data WHERE name = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantName).Scan(&wateredSoilMoisture); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err;
	}

	return wateredSoilMoisture, err
}

func (dbReader DatabaseHandler) GetWateredSoilMoistureFromId(plantId int) (wateredSoilMoisture int, err error) {
	if stmtOut, err = db.Prepare("SELECT watered_soil_moisture FROM plants_data WHERE plant_id = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantId).Scan(&wateredSoilMoisture); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err;
	}

	return wateredSoilMoisture, err
}

func (dbReader DatabaseHandler) getPlantName(plant_id int) (plantName string, err error) {
	if stmtOut, err = db.Prepare("SELECT watered_soil_moisture FROM plants_data WHERE plant_id = ?"); err != nil {
		return "", err
	}

	if err = stmtOut.QueryRow(plant_id).Scan(&plantName); err != nil {
		return "", err
	}

	if err = stmtOut.Close(); err != nil {
		return "", err
	}

	return plantName, err
}

func (dbReader DatabaseHandler) deletePlant(plantName string) (err error) {
	if stmtDel, err = db.Prepare("DELETE FROM plants_data WHERE plantName = ?"); err != nil {
		return err
	}

	if _ ,err = stmtDel.Exec(plantName); err != nil {
		return err
	}

	if err = stmtDel.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) GetNumberOfPlants() (numberOfPlants int, err error) {
	rows, err := db.Query("SELECT COUNT(*) FROM plants_data")

	if err != nil {
		return -1, err
	}

	for rows.Next() {
		err = rows.Scan(&numberOfPlants)
	}

	if err != nil {
		return -1, err
	}

	err = rows.Close()

	if err != nil {
		return -1, err
	}

	return numberOfPlants, nil
}