# Stori Technical Challenge 

_This challenge consist on create a system that processes a file from a mounted directory. The file will contain a list of debit and credit transactions on an account. 
The application should process the file and send summary information to a user in the form of an email._
## Solution

_This solution is developed with GoLang and MySQL on AWS RDS for database and also uses SendGrid email API in order to send the emails.
Also Docker is used to run the service in a single image named "stori_challenge"._

### File Arquitechture
_For this project the structure is the one below_
```
cmd/
├── app
│   ├── config
│   │    └── config.go
│   ├──  database
│   │     └── database.go
│   ├── emailSender
│   │    └── emailSender.go
│   ├── logger
│   │    └── logger.go
│   ├── model
│   │    ├── transaction.go
│   │    └── transactions.go
│   ├── summaryProcessor
│   │   ├── csvProcessor
│   │   │   ├── db.go
│   │   │   └── user.go
│   │   └── summaryReport.go
│   ├── .env
│   └── main.go

```

### Requirements

_You must have installed Docker in order to run the application. You can install it from [here](https://www.docker.com/products/docker-desktop/)._

_You must download the .env file and locate it on /cmd/app/. This file cannot be upload to repository because it has SendGrid API key and it blocks sending emails. Like it contains sensitive information from the database._

###  Installation

_The first step you must do is clone this repository on your own computer. To do that you must have Git installed and then execute the following command:_
```
git clone https://github.com/ErnestoGuevara/StoriChallenge.git
```
_Or you can download the zip file and locate it wherever you want._

### Run
_Once the repository is cloned and Docker already installed. You are going to build the image opening a terminal and being in the same path where the DockerFlie is located, in this case you have to stay in the root path of the repository (/), then you will have to execute the following command:_
```
docker build --tag stori_challenge .
```
_Finally in order to run the image you have to execute the following command, but you have to substitute {myemail@example.com} with your email in order to recive the summary to your email:_
```
docker run -e EMAIL_ADDRESS=myemail@example.com stori_challenge
```

### Results
_Running the image for the first time you will see something like this in your terminal and recive an email with the Summary Report:_
```
DB_INFO: 2023/04/22 04:56:32 logger.go:17: [INFO] ¡Database Connected!
DB_INFO: 2023/04/22 04:56:32 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:33 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:33 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:33 logger.go:17: [INFO] Value inserted on stori_transactions table

--Summary Report--
Total balance is 39.74
Average debit amount: -15.38
Average credit amount: 35.25
Number of transaction in July: 2 
Number of transaction in August: 2 

EMAIL_INFO: 2023/04/22 04:56:34 logger.go:17: [INFO] ¡Email sended!
DB_INFO: 2023/04/22 04:56:34 logger.go:17: [INFO] ¡Database Connected!
DB_INFO: 2023/04/22 04:56:34 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:35 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:35 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:35 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:36 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:36 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:36 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:36 logger.go:17: [INFO] Value inserted on stori_transactions table
DB_INFO: 2023/04/22 04:56:37 logger.go:17: [INFO] Value inserted on stori_transactions table

--Summary Report--
Total balance is 80.50
Average debit amount: -50.00
Average credit amount: 16.31
Number of transaction in August: 1 
Number of transaction in December: 2 
Number of transaction in February: 1 
Number of transaction in March: 1 
Number of transaction in May: 1 
Number of transaction in June: 1 
Number of transaction in September: 1 
Number of transaction in July: 1 

EMAIL_INFO: 2023/04/22 04:56:37 logger.go:17: [INFO] ¡Email sended!
```

_Running the image a second time onwards you will see something like this in your terminal and recive an email with the Summary Report:_
```
DB_INFO: 2023/04/22 04:57:08 logger.go:17: [INFO] ¡Database Connected!

--Summary Report--
Total balance is 39.74
Average debit amount: -15.38
Average credit amount: 35.25
Number of transaction in July: 2 
Number of transaction in August: 2 

EMAIL_INFO: 2023/04/22 04:57:09 logger.go:17: [INFO] ¡Email sended!
DB_INFO: 2023/04/22 04:57:09 logger.go:17: [INFO] ¡Database Connected!

--Summary Report--
Total balance is 80.50
Average debit amount: -50.00
Average credit amount: 16.31
Number of transaction in September: 1 
Number of transaction in July: 1 
Number of transaction in August: 1 
Number of transaction in December: 2 
Number of transaction in February: 1 
Number of transaction in March: 1 
Number of transaction in May: 1 
Number of transaction in June: 1 

EMAIL_INFO: 2023/04/22 04:57:11 logger.go:17: [INFO] ¡Email sended!
```
_And you should receive an email with the following format_

![](https://storichallenge.s3.amazonaws.com/results/Email.jpeg)

### Constraints
_The CSV file must have the Id,Date,Transactions columns._

_Have the .env file located on /cmd/app/ path_
## AWS Implementation
_For the implementation of AWS, the Lambda, S3 and ECR tools were used.
ECR was used to be able to upload the generated image to a repository called "storichallenge" to be able to use it in a container.
Once the image was uploaded to ECR, a Lambda function was created with the option to create it with this image and this facilitate the creation of this function.
As such, a function was developed with Lambda which generates the Summary Report, saves the information in the RDS database and sends the email with the information extracted from a csv format file that is uploaded to the "archvisoCSV/" folder in a bucket named "storichallenge". An event was created in order to the Lambda function is activated as long as a csv file is uploaded to the bucket._

_The flowchart of the aforementioned implementation is represented by the image below._

![](https://storichallenge.s3.amazonaws.com/results/FlowChart.png)

_The Lambda function with the event of uploading a csv file to the bucket is seen as follows._

![](https://storichallenge.s3.amazonaws.com/results/Lambda.jpeg)

_The CloudWatch tool was also used to be able to control the logs, as can be seen in the following image, a successful log after uploading the "client3.csv" file to the bucket._

![](https://storichallenge.s3.amazonaws.com/results/CloudWatch.jpeg)


## Author
Ernesto Ibhar Guevara Gómez
