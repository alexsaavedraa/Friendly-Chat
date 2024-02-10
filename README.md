Hello,
https://dbdiagram.io/d/65c6b006ac844320aed729dd


Postgres Notes:
When you install, make sure you pay attention to what your username and password are
to log into postgres ctl, 
`psql -U username -h localhost`

to see list of databases available, 
`\l` or `\list`

to select a specific database, 
`\c dbname`

to see the tables in the database
`\dt`


In order to run the go backend:
first set your database username and password. on linux/mac use:
`export DB_USERNAME=postgres`
`export DB_PASSWORD=password`
on Windows, use
 $env:DB_USERNAME="postgres"
 $env:DB_PASSWORD="password"


