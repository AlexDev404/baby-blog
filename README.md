# baby-blog

This project just serves to satisfy the requirements of Assignment 1 to _n_ in the [CMPS3162 - Advanced Databases](https://doit.ub.edu.bz/course/view.php?id=70) course

## What this project is

This is a blog and is intended to act as a journal for documenting what us students have learned for each week in the course.

## Prerequisites

You are required to have the following dependencies installed on your system before running this application:

- [PostgreSQL](https://www.postgresql.org/download/)
- [Go](https://go.dev/)
- [Golang-Migrate](https://github.com/golang-migrate/migrate)

## To run

> [!IMPORTANT]
> On First Run (If not, skip to below)
> Open a command-line on either Windows or Linux inside of the project root and execute the following
>
> ```shell
> # Install all dependencies and scaffolding required
> $ make prepare
> $ make initialize

To run the project, you can execute the following each time

```shell
$ make
"Now listening on port http://127.0.0.1:4000"
```


## What it looks like

|Explanation|Preview|
|:-----------|:-------:|
|Provided is a simple user interface for navigating through the blog posts|![image](https://github.com/user-attachments/assets/bb9a1c78-79b8-48e3-b8a9-72ddd771a646)|
|And viewing them as well|![image](https://github.com/user-attachments/assets/25650068-34b9-4ca5-88df-a78673d9f41d)|

There is nothing more left to show. This is a very simple project.

## Technologies used

This project makes use of Go's `net/http` while statically serving the stylesheets.
