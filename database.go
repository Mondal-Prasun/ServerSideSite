package main

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DbCfg struct {
	udb *sql.DB
	cdb *sql.DB
}

type User struct {
	ID        uuid.UUID `json:"uid"`
	Phone     string    `json:"phone"`
	Name      string    `json:"fullName"`
	Gender    string    `json:"gender"`
	BloodType string    `json:"bloodType"`
	CampId    string    `json:"campId"`
	Weight    int       `json:"weight"`
}

type camp struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func initDatabase() (userDb *sql.DB, campDb *sql.DB) {

	cdb, err := sql.Open("sqlite3", "camp.db")

	if err != nil {
		panic("camp database cannot be initilized..." + err.Error())
	}

	campTable := `
	CREATE TABLE IF NOT EXISTS camp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT NOT NULL
);`

	_, err = cdb.Exec(campTable)

	if err != nil {
		panic("table is not crated" + err.Error())
	}

	// dummyCamp := []camp{
	// 	{
	// 		Name:     "Aiims",
	// 		Location: "Kalyani",
	// 	},
	// 	{
	// 		Name:     "Mercy Hospital",
	// 		Location: "kolkata",
	// 	},
	// 	{
	// 		Name:     "Sanjeevni Hospital",
	// 		Location: "kolkata",
	// 	},
	// 	{
	// 		Name:     "SSKM Hospital",
	// 		Location: "Kolkata",
	// 	},
	// }

	// insertCamp := `INSERT INTO camp (name, location) VALUES (?, ?);`

	// for _, camp := range dummyCamp {
	// 	_, err = cdb.Exec(insertCamp, camp.Name, camp.Location)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// }

	udb, err := sql.Open("sqlite3", "donor.db")

	if err != nil {
		panic("user database cannot be initilized..." + err.Error())
	}

	userTable := `
	CREATE TABLE IF NOT EXISTS donor (
    id TEXT PRIMARY KEY,
    phone_number TEXT NOT NULL,
    full_name TEXT NOT NULL,
    gender TEXT CHECK (gender IN ('Male', 'Female', 'Other')) NOT NULL,
    blood_type TEXT CHECK (blood_type IN ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-')) NOT NULL,
    camp_id INTEGER,
    weight REAL CHECK (weight > 0),
    FOREIGN KEY (camp_id) REFERENCES camp(id)
);`

	_, err = udb.Exec(userTable)

	if err != nil {
		panic("table is not crated" + err.Error())
	}

	return udb, cdb
}

func (dbCfg *DbCfg) insertUserTableDonor(user *User) (success bool, error error) {
	insertUser := `INSERT INTO donor (id, phone_number, full_name, gender, blood_type, camp_id, weight)
VALUES (?,?,?,?,?,?,?);
`

	_, err := dbCfg.udb.Exec(insertUser,
		user.ID, user.Phone, user.Name, user.Gender, user.BloodType, user.CampId, user.Weight)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (dbCfg *DbCfg) getAllTableCamp() (interface{}, error) {

	row, err := dbCfg.cdb.Query(`SELECT * FROM camp`)

	if err != nil {
		return nil, err
	}

	var id int
	var name, location string

	type campData struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	var allData []campData

	for row.Next() {
		row.Scan(&id, &name, &location)
		allData = append(allData, campData{
			Id:       id,
			Name:     name,
			Location: location,
		})
	}

	return allData, nil

}
