# footy-predictor

Football predictor backend API

## Code Conventions:

### Database tables :

All table names will be pluralised e.g., **countries** not ~~country~~

### Endpoints CRUD functions names:
Name   | Description
-------| -------------------------------
Create | Insert a new record
Update | Update an existing record
Find   | Retrieve a single record by ID
List   | Retrieve multiple records
Delete | Remove a record

### Domains
* All domains package names will be pluralised e.g., **countries** not ~~country~~
* Every domain package should have 4 files
  
Name                  | Description
--------------------- | -------------------------------
{domain}s_dao.go      | Code for interactions with the database
{domain}s_dao_test.go | Tests for the above code
{domain}s_dto.go      | Types definitions
{domain}s_queries.go  | Query strings constants


