# ABS App

This is the backend side of the ABS App made with Go (Fiber, GORM PostgreSQL) for
ordering and managing at ABS.

This project uses the GNU GPL v2.0 license.

## Development

I developed this app on Ubuntu 22.04 in WSL2, so all the commands you will see
will be Unix commands. It will be slightly different on Windows.

Clone this repo and then `cd` into it. To develop this app locally first you must
have go installed on your machine. If you already have Go, install all dependencies
with the command `go get .`.

### Database

This project uses PostgreSQL. Install it locally on your machine, database migrations
will run automatically when you run the command `go run .`. Sample data is coming soon

### Todo

- ~~Implement auth~~
- implement ordering
  - for anonymous user (e.g. most people)
  - ~~for authenticated users~~
- make a frontend application (maybe as a separate project)
- make `.sql` file for database migrations
- make API documentation, preferrably a dedicated site
- make this api use docker lmao probably don't

After all setup is done, you can go ahead an run `go run .` to run the API.

## How to use

For the end user, are two main endpoints on each you can do CRUD operations:

- `/api/menu`
- `/api/orders`

Each item in `/menu` also has `/variant-values` which contains each menu's variants
and respective prices.

There are also three endpoints for the admin:

- `/admin/menu`
- `/admin/users`
- `/admin/orders`
