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

After all setup is done, you can go ahead and run `go run .`, and Fiber will tell
you that the server is running on `localhost:8080`.

### Database

This project uses PostgreSQL. Install it locally on your machine, database migrations
will run automatically when you run the command `go run .`. For now, the seeder only
has the user roles and one admin user.

Complete sample data is coming soon

### Todo

In order of priority.

- ~~Implement auth~~
  - implement refresh token
- implement ordering
  - for anonymous user (e.g. most people)
  - ~~for authenticated users~~
  - ~~implement order completion~~
  - omit password field in orders
- ~~implement limiting~~
- make `.sql` file for database migrations
  - insert menu data
- make API documentation, preferrably a dedicated site
- make a frontend application, probably with Vue again or Svelte (as a monorepo)
- migrate this API to use Docker (lmao I will probably not do this)
- final configurations and deploy!!

## How to use

For the end user/application, are two main endpoints on each you can do CRUD operations:

- `/api/menu`
- `/api/orders`

Each item in `/menu` also has `/variant-values` which contains each menu's variants
and respective prices.

There are also three endpoints for the admin:

- `/admin/menu`
- `/admin/users`
- `/admin/orders`
