CREATE TABLE IF NOT EXISTS "vacancies"(
    "id" SERIAL PRIMARY KEY,
    "position" TEXT NOT NULL,
    "salary" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS "branches_vacancies"(
    "branch_id" INTEGER REFERENCES branches(id),
    "vacancy_id" INTEGER REFERENCES vacancies(id)
);

