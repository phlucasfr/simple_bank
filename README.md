# Simple Bank
Simple Bank is a digital finance application tailored for Common Users and Merchants. Each user type has a wallet for facilitating transfers. This repository focuses on the core feature: peer-to-peer transfer functionality.

## Stacks
- SQLC: https://sqlc.dev/
- Docker: https://www.docker.com/
- Testify: https://github.com/stretchr/testify
- Postgres: https://www.postgresql.org/
- Migrate: https://github.com/golang-migrate/migrate
- pq: https://github.com/lib/pq

## API Documentation
https://documenter.getpostman.com/view/29417482/2s9Y5cuLTS

<div style="text-align:center;">
  <img src="https://github.com/phlucasfr/simple_bank/blob/main/assets/images/PicPay%20Simplificado.png" style="width: 75%; height: 75%;">
</div>

## **Setup and Usage Guide for the Simple Bank Project**

This project utilizes Docker for setting up the Postgres database and `sqlc` & `migrate` for handling the database operations. Here's a comprehensive guide to using the provided Makefile:

---

### 1. PostgreSQL Setup:

**Command:** `make postgres`

- This command pulls and runs the Postgres v12 Alpine image in a Docker container.
- Container Name: `postgres12`
- Ports: Exposed at `5432` (both internal and host).
- User: `root`
- Password: `secret`

**Example:**
```bash
$ make postgres
```

### 2. Create Database:

**Command:** `make createdb`

This command creates a new PostgreSQL database called picpay_simplificado.
- Uses the username: root and sets the owner as root.

**Example:**

```bash
$ make createdb
```

### 3. Drop Database:

**Command:** `make dropdb`

Deletes the picpay_simplificado database.

**Example:**

```bash
$ make dropdb
```
### 4. Migrate Database (Up):

**Command:** `make migrateup`

Uses the migrate tool to apply all up migrations from the db/migration path to the picpay_simplificado database.
- Uses PostgreSQL connection string with username root, password secret, and no SSL.

**Example:**

```bash
$ make migrateup
```
### 5. Migrate Database (Down):

**Command:** `make migratedown`

Uses the migrate tool to roll back all migrations from the db/migration path in the picpay_simplificado database.

**Example:**

```bash
$ make migratedown
```
### 6. Generate SQL Client Code:

**Command:** `make sqlc`

Uses the sqlc tool to generate client code for the database.

**Example:**

```bash
$ make sqlc
```

### 7. Run Tests:

**Command:** `make test`

Runs all Go tests in the project with verbosity and coverage report.

**Example:**

```bash
$ make test
```

## Contact
If you have any questions or suggestions, please contact the developer at phlucasfr@gmail.com
