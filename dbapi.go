package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	partitionKey = "Email"
)

func tableStruct(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(partitionKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(partitionKey),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest), // ondemand
		TableName:   aws.String(tableName),
	}
}

func validTable(db dynamodb.DynamoDB, tableName string) bool {
	input := &dynamodb.ListTablesInput{}
	isTableName := false
	// 한번에 최대 100개의 테이블만 가지고 올 수 있다.
	// 한 리전에 최대 256개의 테이블이 존재할 수 있다.
	// https://docs.aws.amazon.com/ko_kr/amazondynamodb/latest/developerguide/Limits.html
	for {
		result, err := db.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Fprintf(os.Stderr, "%s %s\n", dynamodb.ErrCodeInternalServerError, err.Error())
				default:
					fmt.Fprintf(os.Stderr, "%s\n", aerr.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			return false
		}

		for _, n := range result.TableNames {
			if *n == tableName {
				isTableName = true
				break
			}
		}
		if isTableName {
			break
		}
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
	return isTableName
}

func hasItem(db dynamodb.DynamoDB, tableName string, primarykey string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(primarykey),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

// AddUser 는 사용자를 추가하는 함수이다.
func AddUser(db dynamodb.DynamoDB) error {
	hasBool, err := hasItem(db, *flagTable, *flagEmail)
	if err != nil {
		return err
	}
	if hasBool {
		return errors.New("The data already exists. Can not add data")
	}
	item := User{
		Email:        *flagEmail,
		UpdateDate:   *flagUpdateDate,
		NameKor:      *flagNameKor,
		NameEng:      *flagNameEng,
		Jobcode:      *flagJobcode,
		Bank:         *flagBank,
		BankAccount:  *flagBankAccount,
		SharesNum:    *flagSharesNum,
		CostHourly:   *flagCostHourly,
		CostWeekly:   *flagCostWeekly,
		CostMonthly:  *flagCostMonthly,
		CostYearly:   *flagCostYearly,
		MonetaryUnit: *flagMonetaryUnit,
		Working:      *flagWorking,
		Projects:     strings.Split(*flagProjects, ","),
	}

	dynamodbJSON, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// SetUser 는 유저자료구조를 수정하는 함수이다.
func SetUser(db dynamodb.DynamoDB) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(*flagTable),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(*flagEmail),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return err
	}
	u := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &u)
	if err != nil {
		return err
	}
	if *flagNameKor != "" && u.NameKor != *flagNameKor {
		u.NameKor = *flagNameKor
	}
	if *flagNameEng != "" && u.NameEng != *flagNameEng {
		u.NameEng = *flagNameEng
	}
	if *flagJobcode != 0 && u.Jobcode != *flagJobcode {
		u.Jobcode = *flagJobcode
	}
	if *flagBank != "" && u.Bank != *flagBank {
		u.Bank = *flagBank
	}
	if *flagBankAccount != "" && u.BankAccount != *flagBankAccount {
		u.BankAccount = *flagBankAccount
	}
	if *flagSharesNum != 0 && u.SharesNum != *flagSharesNum {
		u.SharesNum = *flagSharesNum
	}
	if *flagCostHourly != 0 && u.CostHourly != *flagCostHourly {
		u.CostHourly = *flagCostHourly
	}
	if *flagCostWeekly != 0 && u.CostWeekly != *flagCostWeekly {
		u.CostWeekly = *flagCostWeekly
	}
	if *flagCostMonthly != 0 && u.CostMonthly != *flagCostMonthly {
		u.CostMonthly = *flagCostMonthly
	}
	if *flagCostYearly != 0 && u.CostYearly != *flagCostYearly {
		u.CostYearly = *flagCostYearly
	}
	if *flagMonetaryUnit != "KRW" && u.MonetaryUnit != *flagMonetaryUnit {
		u.MonetaryUnit = *flagMonetaryUnit
	}
	if u.Working != *flagWorking {
		u.Working = *flagWorking
	}
	if *flagProjects != "" && reflect.DeepEqual(u.Projects, strings.Split(*flagProjects, ",")) == false {
		u.Projects = strings.Split(*flagProjects, ",")
	}
	u.UpdateDate = *flagUpdateDate
	dynamodbJSON, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}
	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// RmUser 는 유저자료구조를 사용자를 삭제하는 함수이다.
func RmUser(db dynamodb.DynamoDB) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(*flagTable),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(*flagEmail),
			},
		},
	}
	_, err := db.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers 는 사용자를 가지고오는 함수이다.
func GetUsers(db dynamodb.DynamoDB, word string) error {
	proj := expression.NamesList(
		expression.Name("Email"),
		expression.Name("NameKor"),
		expression.Name("NameEng"),
		expression.Name("Jobcode"),
		expression.Name("Bank"),
		expression.Name("BankAccount"),
		expression.Name("SharesNum"),
		expression.Name("CostHourly"),
		expression.Name("CostWeekly"),
		expression.Name("CostMonthly"),
		expression.Name("CostYearly"),
		expression.Name("MonetaryUnit"),
		expression.Name("Working"),
		expression.Name("Projects"),
	)
	filt1 := expression.Name("Email").Contains(word)
	filt2 := expression.Name("NameKor").Contains(word)
	filt3 := expression.Name("NameEng").Contains(word)
	filt4 := expression.Name("Projects").Contains(word)

	expr, err := expression.NewBuilder().
		WithFilter(filt1).
		WithFilter(filt2).
		WithFilter(filt3).
		WithFilter(filt4).
		WithProjection(proj).
		Build()
	if err != nil {
		return err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(*flagTable),
	}
	result, err := db.Scan(params)
	if err != nil {
		return err
	}
	for _, i := range result.Items {
		u := User{}
		err = dynamodbattribute.UnmarshalMap(i, &u)
		if err != nil {
			return err
		}
		fmt.Println(u)
	}
	return nil
}
