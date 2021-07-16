# Golang CRUD test task

This application is a Golang api structure with gorm and fully version manageable. It provides options to improve the
performance of the developer.

## 🙇 Application Requirement

1. Ensure Golang is an installed or not, if installed then make sure its version is ```>= 1.13``` and if not installed
   then visit https://golang.org/dl/
    * To check the version of Golang use this

```bash   
go version
```

2. Make sure GOPATH is a set or not, if not visit https://github.com/golang/go/wiki/SettingGOPATH

3. Need ```MySQL``` as a Database.

### 🎵 Note

- To send an email using SMTP please enable the ```Less Secure App Access``` setting.
- Please create a blank database in your ```Mysql``` and add that database name in ```ConnectionString``` variable
  of ```.env``` file as mention in step 2

## 🛠️ Start the application locally

1. Clone the repository in Golang's ```src``` folder

2. Change the below ```ENV variables``` values as per your system environment in ```.env``` file.

```
ConnectionString = {username}:{password}@tcp({host}:{port})/{database_name}?charset=utf8&parseTime=True&loc=Local
Host = smtp.xxxxxxxx.com
Sender = xxxxx@xxxx.com
Password = xxxxxxxxxx
```

3. Run `go mod vendor` to install all the required dependency
4. From project ```root``` directory, open a terminal and run

```bash 
go run main.go //to start the process and if its first time then migrate all the required tables
```

then after to run the seeder execute below command

```bash
go run main.go seed  //this will run all the required seeders
```

5. Now the system is ready for the execution of APIs. Please import the below POSTMAN collection URL in your postman.

```   
https://www.getpostman.com/collections/f3715a1cb5227271a567
```

6. Below are the list of APIs available in the POSTMAN collection

```
1. Signup
2. Verify Email
3. Login
4. Forgot Password
5. Reset Password
6. List
7. Delete
```

🌟 You are all set!
