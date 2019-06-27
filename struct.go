package main

import (
	"fmt"
	"strings"
)

// User 는 사용자 정보를 다루는 자료구조이다.
type User struct {
	Email        string   // Lazypic 이메일. partition key
	UpdateDate   string   // 업데이트 날짜. sort key
	NameKor      string   // 이름
	NameEng      string   // 영문이름
	Jobcode      int      // 업종코드
	Bank         string   // 은행명
	BankAccount  string   // 계좌번호
	SharesNum    int64    // 주식량(지분)
	CostHourly   int64    // 시급
	CostWeekly   int64    // 주급
	CostMonthly  int64    // 월급
	CostYearly   int64    // 연봉
	MonetaryUnit string   // 단위. KRW
	Working      bool     // 상태
	Projects     []string // 참여 프로젝트
}

// Lazypic 에서 프리랜서로 활용할 때 견적서에 등록해야할 업종코드는 아래와 같다.
// 아래 코드는 Jobcode 로 사용된다.
//
// 저술가,시나리오,작가: 940100
// 회화,만화가,삽화가: 940200
// 작곡가,작사가: 940301
// 촬영보조: 940500
// 프로그래머: 940909
// 개인과외,서비스: 940903
// 자문,지도료,고문료: 940600

func (u User) String() string {
	return fmt.Sprintf(`
ID: %s Working: %t
Name: %s(%s)
Shares: %d
Jobcode: %d
Bank: %s %s
Cost: h:%d w:%d m:%d y:%d %s
Projects: %s`,
		u.Email,
		u.Working,
		u.NameKor,
		u.NameEng,
		u.SharesNum,
		u.Jobcode,
		u.Bank,
		u.BankAccount,
		u.CostHourly,
		u.CostWeekly,
		u.CostMonthly,
		u.CostYearly,
		u.MonetaryUnit,
		strings.Join(u.Projects, ","),
	)
}
