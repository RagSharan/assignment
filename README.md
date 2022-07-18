# assignment

#Run Application By Command in vs code

go run cmd/svr/main.go

Run through build file
./main.exe

# Answer APIs
1- Post API- http://localhost:8080/api/v1/answer
        Body- {
              "Key": "Q1",
              "Value": 23
        }
        
2- Get API-  
        GetById- http://localhost:8080/api/v1/answer/62d4558d2000fc95e9f117ff
        GetByOtherParams- http://localhost:8080/api/v1/answer?id=62d4546db78721f67bf825de

3- Put API- http://localhost:8080/api/v1/answer
        Body   {    "id": "62d4558f2000fc95e9f11803",
                    "key": "Q1",
                    "value": "answer"
              }

4- Delete API-  DeleteById
                http://localhost:8080/api/v1/answer/62d4558f2000fc95e9f11803

# Events API



# DB Mongodb docker container url
url- mongodb://localhost:27017

All the code is inside pkg folder
