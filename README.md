# ABS App

This is the backend side of the ABS App made with Go (Fiber, GORM PostgreSQL) for
ordering and managing at ABS.

This project uses the GNU GPL v2.0 license.

## Table of contents

- [Development](#development)
  - [Database](#database)
  - [Todo](#database)
- [How to use](#how-to-use)
  - [Auth](#auth)
    - [Sign up](#sign-up)
    - [Log in](#log-in)
    - [Log out](#log-out)
  - [Menu](#menu)
    - [Get all menu items](#get-all-menu-items)
    - [Get menu item by ID](#get-menu-item-by-id)
  - [Orders](#orders)
    - [Get orders](#get-orders)
    - [Get order by ID](#get-order-by-id)
    - [Place a new order](#place-a-new-order)
- [Admin](#admin)
  - [Menu administration](#menu-administration)
    - [Create new menu item](#create-new-menu-item)
    - [Update existing menu item](#update-an-existing-menu-item)
    - [Delete a menu item](#delete-a-menu-item)
    - [Insert new price for a menu item](#insert-new-price-for-a-menu-item)
    - [Update existing price for a menu item](#update-existing-price-of-a-menu-item)
    - [Delete prices of a menu item](#delete-prices-of-a-menu-item)
  - [Orders administration](#orders-administration)
    - [Get all orders](#get-all-orders)
    - [Complete an order](#complete-an-order)
  - [Users administration](#users-administration)
    - [Get all users](#get-all-users)
    - [Get user by ID](#get-user)

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
- ~~implement ordering~~
  - ~~for anonymous user (e.g. most people)~~
  - ~~for authenticated users~~
  - ~~implement order completion~~
  - omit password field in orders
- ~~implement limiting~~
- implement better error handling for various edge cases
  - ~~signup~~
  - login
- ~~make `.sql` file for database migrations~~
- make API documentation (in progress)
- migrate this API to use Docker (lmao I will probably not do this)
- final configurations and deploy!!
- make a frontend application probably

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

All requests have a data (object or array) and error (bool) property. The `err`
field will be true if there are any errors, client-side or server-side. In that case,
a field named `message` will tell you exactly what happened.

```json
{
  "err": true,
  "message": "..."
}
```

For example, if you tried to access a protected endpoint (say `/api/orders`)
without valid credentials, you would get the following response:

```json
{
  "err": true,
  "message": "cannot verify token"
}
```

An exhaustive list of error messages will be documented soon.

### Auth

#### Sign up

To sign up a user, issue a POST request to `/signup` with a JSON body with the
following fields:

- `name` (string),
- `email`(string): has to be valid email
- `password`(string): has to be >= 8 characters long

```json
{
  "name": "Muhammad Rava",
  "email": "rava@gmail.com",
  "password": "secret-password"
}
```

A successful response looks like this:

```json
{
  "err": false,
  "message": "signup success! redirect user to login page"
}
```

You will get the following error if email is not valid or password is less then 8
characters long:

```json
{
  "err": true,
  "message": "email has to be valid, and password has to be at least 8 characters long"
}
```

#### Log in

To log a user in, issue a POST request to `/signin` with a JSON body with fields `email` and
`password`:

```json
{
  "email": "rava@gmail.com",
  "password": "secret-password"
}
```

A successful response looks like:

```json
{
  "err": false,
  "message": "successfully logged in"
}
```

The API will issue a cookie for the client to log in.

#### Log out

To log a user out, issue a POST request to `/logout` with no body. A successful
request yields a response like this:

```json
{
  "err": false,
  "message": "successfully logged out"
}
```

### Menu

#### Get all menu items

To get all menu items, you can issue a GET request to `/api/menu`. You can also search
by name or filter by type by adding an URL parameter like so: `/api/menu?name=searchterm&type_id=10`.

A successful response looks like this:

<details>
    <summary>
        successful response example
    </summary>

```json
{
  "data": [
    {
      "id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
      "name": "Espresso",
      "type_id": 1,
      "type": {
        "id": 1,
        "type": "kopi"
      },
      "created_at": "2024-01-19T08:56:35.169389+07:00",
      "updated_at": "2024-01-19T08:56:35.169389+07:00",
      "variant_values": [
        {
          "menu_id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
          "option_id": 1,
          "option_value_id": 2,
          "option": {
            "id": 1,
            "name": "temp"
          },
          "option_value": {
            "id": 2,
            "option_id": 1,
            "value": "hot"
          },
          "price": 6000
        },
        {
          "menu_id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
          "option_id": 1,
          "option_value_id": 1,
          "option": {
            "id": 1,
            "name": "temp"
          },
          "option_value": {
            "id": 1,
            "option_id": 1,
            "value": "iced"
          },
          "price": 10000
        }
      ]
    }
  ],
  "error": false
}
```

</details>

#### Get menu item by ID

You can also get a menu item by ID by issuing a GET request to `/api/menu/valid-menu-id`
Where `valid-menu-id` is a valid menu UUID. A successful response looks like this:

<details>
    <summary>
        Successful response example
    </summary>

```json
{
  "data": {
    "id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
    "name": "Espresso",
    "type_id": 1,
    "type": {
      "id": 1,
      "type": "kopi"
    },
    "created_at": "2024-01-19T08:56:35.169389+07:00",
    "updated_at": "2024-01-19T08:56:35.169389+07:00",
    "variant_values": [
      {
        "menu_id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
        "option_id": 1,
        "option_value_id": 2,
        "option": {
          "id": 1,
          "name": "temp"
        },
        "option_value": {
          "id": 2,
          "option_id": 1,
          "value": "hot"
        },
        "price": 6000
      },
      {
        "menu_id": "b983ff25-c532-43ea-aee7-fc75eaa7c2bb",
        "option_id": 1,
        "option_value_id": 1,
        "option": {
          "id": 1,
          "name": "temp"
        },
        "option_value": {
          "id": 1,
          "option_id": 1,
          "value": "iced"
        },
        "price": 10000
      }
    ]
  },
  "err": false
}
```

</details>

### Orders

#### Get orders

To get orders you can issue a GET request to `/api/orders`; this gets all
**complete** orders for the current users by default. You can filter by status
by adding the query parameter `is_completed` to the URL, for example to get incomplete
orders:

`/api/orders?is_completed=false`

If successful, you will get a response following this structure:

<details>
    <summary>
        JSON response
    </summary>

```json
{
  "data": [
    {
      "id": "686156e2-993b-4518-ab42-e757335fcd75",
      "user_id": "fad6c002-cbba-48dc-81d8-d56a17f5428c",
      "user": {
        "id": "fad6c002-cbba-48dc-81d8-d56a17f5428c",
        "name": "Abiman You",
        "email": "abim@abim.com",
        "password": "$2a$14$DhwQKFBX8ZVoIJqBokT6guiKeJ063uBpekz.ZM1ISncnhe.xm1/Qe",
        "role_id": 1,
        "role": {
          "id": 1,
          "name": "user"
        },
        "created_at": "2024-02-19T10:01:05.918176+07:00",
        "updated_at": "2024-02-19T10:01:05.918176+07:00"
      },
      "created_at": "2024-02-20T15:35:39.374745+07:00",
      "is_completed": false,
      "completed_at": "0001-01-01T00:00:00Z",
      "order_details": [
        {
          "order_id": "686156e2-993b-4518-ab42-e757335fcd75",
          "menu_id": "ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67",
          "menu_name": "Banana Strawberry",
          "menu_type": "pizza",
          "menu_option_id": 35,
          "menu_option": "pizza topping",
          "menu_option_value_id": 36,
          "menu_option_value": "regular",
          "quantity": 2,
          "total_price": 90000
        }
      ]
    }
  ],
  "err": false
}
```

</details>

#### Get order By ID

You can also get individual orders by issuing a GET request to `/api/orders/:id`
where `id` is a valid order ID.

#### Place a new order

This endpoint is protected by auth. Makes sure you have logged in before making
a request. To place a new order, issue a POST request with the following body:

```json
{
  "order_details": [
    {
      "menu_id": "valid menu ID here",
      "menu_option_id": 999,
      "menu_option_value_id": 999,
      "quantity": 1
    }
  ]
}
```

(UNSTABLE: anonymous ordering) If placing a new order for an anonymous user, you can supply a username:

```json
{
  "username": "Rava Basya",
  "order_details": [
    {
      "menu_id": "valid menu ID here",
      "menu_option_id": 999,
      "menu_option_value_id": 999,
      "quantity": 1
    }
  ]
}
```

The `order_details` field is an array containing however many items you want to
order. The item is an object with the following fields:

- `menu_id` (uuid): ID of menu item
- `menu_option_id` (int): ID of menu option
- `menu_option_value_id` (int): Id of menu option value

You need to specify `menu_option_id` and `menu_option_value_id` according to the
menu type. For example, the pizza menu type has `35` for `option_id` (pizza topping)
and `36` (regular) for one of the `option_value_id`, therefore your order should
look like this:

```json
{
  "order_details": [
    {
      "menu_id": "bf5aad5a-4d81-4ab2-ae64-d0dad9b77061",
      "menu_option_id": 35,
      "menu_option_value_id": 36,
      "quantity": 1
    }
  ]
}
```

You can look at the available options and the corresponding option values when
getting menu information.

## Admin

All admin endpoints are located in `/admin`. Here you can administer orders, users
and menu items as admin, but to do those things you need to be authenticated.
Head over to `/signin` and sign the admin in.

### Menu administration

#### Create new menu item

to create a new menu item, issue a POST request to `/admin/menu` with the following
body:

```json
{
  "name": "a new delicious drink",
  "type_id": 6,
  "variant_values": [
    {
      "option_id": 1,
      "option_value_id": 37,
      "price": 8000
    },
    {
      "option_id": 1,
      "option_value_id": 2,
      "price": 10000
    }
  ]
}
```

#### Update an existing menu item

TO update an existing menu item, issue a POST request to `/admin/menu` with the following
body:

```json
{
  "name": "Minuman ngetes doang",
  "variant_values": [
    {
      "option_id": 1,
      "option_value_id": 2,
      "price": 10000
    }
  ]
}
```

The server will only update non nil values. In the example above, the server will
update name and the price of menu item where `option_id` is 1 and `option_value_id`is 2.

#### Delete a menu item

To delete a menu item, issue a DELETE request to `/admin/menu/:id` where `id`
is a valid menu item ID of type UUID you want to delete. A successful response looks
like this:

```json
{
  "err": false,
  "message": "menu item successfully deleted"
}
```

You can also batch delete (make a DELETE request to `/admin/menu/`) by supplying
the UUIDs of menu items in an array like so:

```json
["menu-id-1", "menu-id-2"]
```

#### Insert new price for a menu item

To insert a new menu price of a menu item, issue a POST request to `/admin/menu/:id/variant_values`
where `id` is the menu item ID of type UUID.

#### Update existing price of a menu item

To edit a specific price of a menu item, issue a PATCH `/admin/menu/:id/variant_values`
where `id` is a valid menu item ID of type UUID. Attach a body specifying which
combination of option_id and option_value_id you want to edit, and the new price:

```json
[
  {
    "option_id": 1,
    "option_value_id": 1,
    "price": 15000
  }
]
```

#### Delete prices of a menu item

To delete a price of a menu item, issue a DELETE request to `/admin/menu/:id/variant_values`
where `id` is a valid menu item ID of type UUID with the following body:

```json
{
  "": "blom bentar ya"
}
```

### Orders administration

#### Get all orders

To get all orders of all users, issue a GET request to `/admin/orders/`. You can
also get orders by ID (`/admin/orders/:id`) where id is a valid order ID.

<details>
    <summary>Successful response example</summary>

```json
{
  "data": [
    {
      "id": "26a33223-da03-4a5c-8bbb-f7fd6944abef",
      "user_id": "46e11084-0baa-4d2f-bf52-6f6a93a78619",
      "user": {
        "id": "46e11084-0baa-4d2f-bf52-6f6a93a78619",
        "name": "Muhammad Rava",
        "email": "rava@gmail.com",
        "password": "password",
        "role_id": 1,
        "role": {
          "id": 1,
          "name": "user"
        },
        "created_at": "2024-02-21T10:12:43.369861+07:00",
        "updated_at": "2024-02-21T10:12:43.369861+07:00"
      },
      "created_at": "2024-02-21T10:18:19.674673+07:00",
      "is_completed": true,
      "completed_at": "2024-02-21T10:26:54.973271+07:00",
      "order_details": [
        {
          "order_id": "26a33223-da03-4a5c-8bbb-f7fd6944abef",
          "menu_id": "ab2a528c-9c5b-45d0-ba7f-ce91a97c6b67",
          "menu_name": "Banana Strawberry",
          "menu_type": "pizza",
          "menu_option_id": 35,
          "menu_option": "pizza topping",
          "menu_option_value_id": 36,
          "menu_option_value": "regular",
          "quantity": 2,
          "total_price": 90000
        }
      ]
    }
  ],
  "err": false
}
```

</details>

<details>
    <summary>error response</summary>

```json
{
  "err": true,
  "message": "error when querying database"
}
```

</details>

#### Complete an order

To complete an order, issue a PATCH request (`/admin/orders/:id`) where `id` is
a valid order ID.

<details>
    <summary>Successful response example</summary>

```json
{
  "err": false,
  "message": "order of ID <id> succesfully completed"
}
```

</details>

### Users administration

#### Get all users

To get all users, issue a GET request to `/admin/users`.

<details>
    <summary>
        Succcesful response example
    </summary>

```json
{
  "data": [
    {
      "id": "fad6c002-cbba-48dc-81d8-d56a17f5428c",
      "name": "Abiman You",
      "email": "abim@abim.com",
      "password": "",
      "role_id": 1,
      "role": {
        "id": 1,
        "name": "user"
      },
      "created_at": "2024-02-19T10:01:05.918176+07:00",
      "updated_at": "2024-02-19T10:01:05.918176+07:00"
    }
  ],
  "err": false
}
```

</details>

#### Get user by ID

You can also get a specific user by adding an id (`/admin/users/:id`).

<details>
    <summary>
        Successful response example
    </summary>

```json
{
  "data": {
    "id": "fad6c002-cbba-48dc-81d8-d56a17f5428c",
    "name": "Abiman You",
    "email": "abim@abim.com",
    "password": "",
    "role_id": 1,
    "role": {
      "id": 1,
      "name": "user"
    },
    "created_at": "2024-02-19T10:01:05.918176+07:00",
    "updated_at": "2024-02-19T10:01:05.918176+07:00"
  },
  "err": false
}
```

</details>
