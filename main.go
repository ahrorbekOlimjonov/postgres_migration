package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Respons struct {
	Stores []*Store
}

type Store struct {
	ID       int
	Name     string
	Branches []*Branch
}
type Branch struct {
	ID          int
	Name        string
	PhoneNumber []string
	StoreID     int
	Vacancies   []*Vacancy
}
type Vacancy struct {
	ID       int
	Position string
	Salary   int
	BranchId int
}

func connect() *sql.DB {

	connect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "ahrorbek", "3108", "market")

	db, err := sql.Open("postgres", connect)

	if err != nil {
		panic(err)
	}
	return db

}

func insertStore(db *sql.DB, stores []Store) {

	for _, s := range stores {

		_, err := db.Exec("INSERT INTO stores(name) VALUES($1)", s.Name)
		if err != nil {
			panic(err)
		}

	}
}

func insertBranch(db *sql.DB, stores []Store) {

	for _, s := range stores {
		var store_id int

		err := db.QueryRow("Select id from stores where name = $1", s.Name).Scan(&store_id)

		if err != nil {
			panic(err)
		}

		for _, b := range s.Branches {

			_, err := db.Exec("INSERT INTO branches(name, phonenumber, store_id) VALUES"+
				"($1, $2, $3)", b.Name, pq.Array(b.PhoneNumber), store_id)

			if err != nil {
				panic(err)
			}
		}
	}

}

func insertVacancy(db *sql.DB, stores []Store) {

	for _, s := range stores {
		for _, b := range s.Branches {
			for _, v := range b.Vacancies {

				_, err := db.Exec("INSERT INTO vacancies(position, salary) VALUES($1, $2)", v.Position, v.Salary)

				if err != nil {
					return
				}

			}
		}
	}
}

func insertBranchesVacancies(db *sql.DB, stores []Store) {
	for _, s := range stores {

		for _, b := range s.Branches {
			var branch_id int

			err := db.QueryRow("Select id from branches where name = $1", b.Name).Scan(&branch_id)

			if err != nil {
				panic(err)
			}
			for _, v := range b.Vacancies {

				var vacancy_id int

				err := db.QueryRow("Select id from vacancies where position = $1", v.Position).Scan(&vacancy_id)

				if err != nil {
					panic(err)
				}

				_, err = db.Exec("Insert into "+
					"branches_vacancies(branch_id, vacancy_id) "+
					"values ($1, $2)",
					branch_id, vacancy_id)
				fmt.Println("Run")

				if err != nil {
					panic(err)
				}

			}
		}
	}
}

func reloadDatabase(db *sql.DB) {

	_, err7 := db.Exec("DROP TABLE branches_vacancies")
	if err7 != nil {
		panic(err7)
	}
	_, err4 := db.Exec("DROP TABLE vacancies")
	if err4 != nil {
		panic(err4)
	}
	_, err2 := db.Exec("DROP TABLE branches")
	if err2 != nil {
		panic(err2)
	}
	_, err := db.Exec("DROP TABLE stores")
	if err != nil {
		panic(err)
	}

	_, err1 := db.Exec("CREATE TABLE " +
		"stores(id serial primary key , name varchar(64))")

	if err1 != nil {
		panic(err1)
	}

	_, err3 := db.Exec("CREATE TABLE " +
		"branches(id serial primary key , name varchar(64), phoneNumber varchar array, store_id integer  references stores(id))")
	if err3 != nil {
		panic(err3)
	}

	_, err5 := db.Exec("CREATE TABLE " +
		"vacancies(id serial primary key , position varchar(64), salary int)")
	if err5 != nil {
		panic(err5)
	}

	_, err6 := db.Exec("CREATE TABLE " +
		"branches_vacancies(branch_id integer references branches(id), vacancy_id integer references vacancies(id))")
	if err6 != nil {
		panic(err6)
	}

}

func deleteData(db *sql.DB) {

	_, err1 := db.Exec("DELETE FROM vacancy WHERE id = ($1)", 8)
	if err1 != nil {
		panic(err1)
	}

}

func updateDta(db *sql.DB) {

	branch := Branch{
		ID:   1,
		Name: "Abay",
	}

	_, err := db.Exec("UPDATE branch SET name = ($1) WHERE id = ($2)", branch.Name, branch.ID)
	if err != nil {
		panic(err)
	}
}

func getStoreData(db *sql.DB) {

	fmt.Println("	Table Store")

	rows, err := db.Query("SELECT id, name FROM store")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

}

func getBranchData(db *sql.DB) {

	fmt.Println("	Table Branch")

	rows, err := db.Query("SELECT id, name, phonenumber, store_id FROM branch")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var phonenumber string
		var store_id int

		if err := rows.Scan(&id, &name, &phonenumber, &store_id); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Name: %s, Phone number: %v, Store_id: %d\n", id, name, phonenumber, store_id)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func getVacancyData(db *sql.DB) {

	fmt.Println("	Table Vacancy")

	rows, err := db.Query("SELECT id, position, salary, branch_id FROM vacancy")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var position string
		var salary int
		var branch_id int

		if err := rows.Scan(&id, &position, &salary, &branch_id); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Position: %s, Salary: %d, Branch_id: %d\n", id, position, salary, branch_id)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func getAllData(db *sql.DB) {
	rows, err := db.Query("select s.Name, b.name, v.position, v.salary from store s " +
		"join branch b on b.store_id = s.id " +
		"join vacancy v on v.branch_id  = b.id")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var storeName, branchName, position string
		var salary float64
		err = rows.Scan(&storeName, &branchName, &position, &salary)
		if err != nil {
			panic(err)
		}
		fmt.Println(storeName, branchName, position, salary)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

func catch_err(stores []Store) {}

func main() {
/*
	
	db := connect()
	defer db.Close()

	stores := []Store{
		{
			Name: "Karzinka",
			Branches: []*Branch{
				{
					Name:        "Kohinur",
					PhoneNumber: []string{"998931234567", "998931234567"},
					Vacancies: []*Vacancy{
						{
							Position: "Driver",
							Salary:   300,
						},
						{
							Position: "Saler",
							Salary:   400,
						},
					},
				},
				{
					Name:        "Beruniy",
					PhoneNumber: []string{"998931234567", "998931234567"},
					Vacancies: []*Vacancy{
						{
							Position: "Guard",
							Salary:   250,
						},
						{
							Position: "Saler",
							Salary:   350,
						},
					},
				},
			},
		},
		{
			Name: "Makro",
			Branches: []*Branch{
				{
					Name:        "Olmazor",
					PhoneNumber: []string{"998931234567", "998931234567"},
					Vacancies: []*Vacancy{
						{
							Position: "Quality_checker",
							Salary:   399,
						},
						{
							Position: "Driver",
							Salary:   270,
						},
					},
				},
				{
					Name:        "Chorsu",
					PhoneNumber: []string{"998931234567", "998931234567"},
					Vacancies: []*Vacancy{
						{
							Position: "Cleaner",
							Salary:   330,
						},
						{
							ID:       8,
							Position: "manager",
							Salary:   500,
						},
					},
				},
			},
		},
	}
	catch_err(stores)

	reloadDatabase(db)
	insertStore(db, stores)
	insertBranch(db, stores)
	insertVacancy(db, stores)
	insertBranchesVacancies(db, stores)
	updateDta(db)
	deleteData(db)
	getStoreData(db)
	getBranchData(db)
	getVacancyData(db)
	getAllData(db)

	resp := Respons{}
	
	sRows, err := db.Query("select id, name from stores")
	if err != nil {
		return
	}
	for sRows.Next() {
	
		store := Store{}
		err := sRows.Scan(
			&store.ID,
			&store.Name,
		)
		if err != nil {
			return
		}
	
		bRows, err := db.Query("SELECT id, name, phonenumber from branches where store_id = $1", store.ID)
		if err != nil {
			return
		}
		for bRows.Next() {
			branch := Branch{}
			err := bRows.Scan(
				&branch.ID,
				&branch.Name,
				pq.Array(&branch.PhoneNumber),
			)
			if err != nil {
				return
			}
	
			vRows, err := db.Query("SELECT v.id, v.position, v,salary "+
				"FROM vacancies v "+
				"JOIN branches_vacancies bv "+
				"ON v.id = bv.vacancy_id "+
				"JOIN branches b "+
				"ON b.id = bv.branch_id"+
				"WHERE b.id = $1", branch.ID)
	
			for vRows.Next() {
				vacancy := Vacancy{}
				err := vRows.Scan(
					&vacancy.ID,
					&vacancy.Position,
					&vacancy.Salary,
				)
				if err != nil {
					return
				}
	
				branch.Vacancies = append(branch.Vacancies, &vacancy)
			}
			store.Branches = append(store.Branches, &branch)
		}
		resp.Stores = append(resp.Stores, &store)
	}
	
	for _, store := range resp.Stores {
		fmt.Println(store)
	}

	
	resp := Respons{}

	storeRows, err := db.Query("SELECT id, name FROM stores")
	if err != nil {
		fmt.Println(err)
		return
	}

	for storeRows.Next() {
		store := Store{}
		err := storeRows.Scan(
			&store.ID,
			&store.Name,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		branchRows, err := db.Query("SELECT id, name, phonenumber from branches WHERE store_id = $1", store.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		for branchRows.Next() {
			branch := Branch{}
			err := branchRows.Scan(
				&branch.ID,
				&branch.Name,
				pq.Array(&branch.PhoneNumber),
			)
			if err != nil {
				fmt.Println(err)
				return
			}

			vacancyRows, err := db.Query("SELECT v.id, v.position, v.salary FROM vacancies v JOIN branches_vacancies br ON v.id = br.vacancy_id JOIN branches b ON b.id = br.branch_id where b.id = $1", branch.ID)
			if err != nil {
				fmt.Println(err)
				return
			}

			for vacancyRows.Next() {
				vacancy := Vacancy{}
				err := vacancyRows.Scan(
					&vacancy.ID,
					&vacancy.Position,
					&vacancy.Salary,
				)
				if err != nil {
					fmt.Println(err)
					return
				}

				branch.Vacancies = append(branch.Vacancies, &vacancy)
			}

			if err != nil {
				fmt.Println(err)
				return
			}

			store.Branches = append(store.Branches, &branch)
		}
		resp.Stores = append(resp.Stores, &store)
	}

	for _, store := range resp.Stores {
		fmt.Println(store)

		for _, branch := range store.Branches {
			fmt.Println(branch)
			for _, v := range branch.Vacancies {
				fmt.Println(v)

			}

		}
	}
	*/



	
}
