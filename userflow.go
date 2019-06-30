package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	now = time.Now()
	// db setting
	flagRegion  = flag.String("region", "ap-northeast-2", "AWS region name")
	flagProfile = flag.String("profile", "lazypic", "AWS Credentials profile name")
	flagTable   = flag.String("table", "userflow", "AWS Dynamodb table name")

	// mode and partition key
	flagAdd = flag.Bool("add", false, "user addition mode")
	flagSet = flag.Bool("set", false, "user update mode")
	flagRm  = flag.Bool("rm", false, "user remove mode")

	// date
	flagHelp       = flag.Bool("help", false, "print help")
	flagUpdateDate = flag.String("updatedate", now.Format(time.RFC3339), "user update date")

	// attributes
	flagNameKor      = flag.String("namekor", "", "korean user name")
	flagNameEng      = flag.String("nameeng", "", "english user name")
	flagEmail        = flag.String("email", "", "lazypic email")
	flagJobcode      = flag.Int("jobcode", 0, "job code number")
	flagBank         = flag.String("bank", "", "bank name")
	flagBankAccount  = flag.String("bankaccount", "", "bank account number")
	flagSharesNum    = flag.Int64("sharenum", 0, "shares number")
	flagCostHourly   = flag.Int64("costhourly", 0, "cost hourly")
	flagCostWeekly   = flag.Int64("costweekly", 0, "cost weekly")
	flagCostMonthly  = flag.Int64("costmonthly", 0, "cost monthly")
	flagCostYearly   = flag.Int64("costyearly", 0, "cost yearly")
	flagMonetaryUnit = flag.String("monetaryunit", "KRW", "monetary unit")
	flagWorking      = flag.Bool("working", false, "is working?")
	flagProjects     = flag.String("projects", "", "projectname")
	flagSearchword   = flag.String("search", "", "search word")
)

func main() {
	log.SetPrefix("userflow: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*flagRegion)},
		Profile:           *flagProfile,
	}))
	db := dynamodb.New(sess)

	// 테이블이 존재하는지 점검하고 없다면 테이블을 생성한다.
	if !validTable(*db, *flagTable) {
		_, err := db.CreateTable(tableStruct(*flagTable))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created table:", *flagTable)
		fmt.Println("Please try again in one minute.")
		os.Exit(0)
	}
	if *flagHelp {
		flag.Usage()
	}
	if *flagAdd && *flagEmail != "" {
		err := AddUser(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagSet && *flagEmail != "" {
		err := SetUser(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagRm && *flagEmail != "" {
		user, err := user.Current()
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		err = RmUser(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		err := GetUsers(*db, *flagSearchword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}
