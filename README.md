![userflow](figures/risalogo.svg)

- Lazypic 유저 매니징툴 입니다.
- Lazypic 내부툴 PlanA에 합병되었습니다. 사용되지 않습니다.

### 사용법
사용자 관리는 고유한 키가되는 Email을 기준으로 관리합니다.

User 추가

```bash
$ userflow -add -namekor 김한웅 -nameeng JasonKim -email woong@lazypic.org -jobcode 940909 -bank 우리은행 -bankaccount 092340913412 -projects circle
```

User 수정

```bash
$ userflow -set -email woong@lazypic.org -projects circle,csi
```

User 삭제

```bash
$ sudo userflow -rm -email woong@lazypic.org
```

User 검색

```bash
$ userflow -searchword [검색어]
```

### 인수

#### DB셋팅

- region: 기본값 "ap-northeast-2", AWS 리전명
- profile: 기본값 "lazypic", AWS Credentials profile 이름
- table: 기본값 "userflow", AWS Dynamodb tablbe 이름

#### 모드
- add: add mode on
- set: set mode on
- rm: rm mode on

#### 속성
- namekor: 한글이름
- nameeng: 영문이름
- email: 이메일
- jobcode: 업종코드
    - 저술가,시나리오,작가: 940100
    - 회화,만화가,삽화가: 940200
    - 작곡가,작사가: 940301
    - 촬영보조: 940500
    - 프로그래머: 940909
    - 개인과외,서비스: 940903
    - 자문,지도료,고문료: 940600
- bank: 은행명
- bankaccount: 계좌번호
- sharenum: 주식수
- costhourly: 시간당임금
- costweekly: 주당임금
- costmonthly: 월급여
- costyearly: 연봉
- monetaryunit", 기본값 "KRW", 화폐단위
- working: 기본값 false, 현재 일하는 상태인지 체크
- projects: 기본값 "", 참여중인 프로젝트
- searchword: 기본값 "", 검색어

#### 기타
- help: 도움말 출력
- updatedate: 사용자 업데이트 날짜를 임의로 변경시 사용

### AWS DB권한 설정
AWS DB접근 권한을 설정할 계정에 아래 권한을 부여합니다.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ListAndDescribe",
            "Effect": "Allow",
            "Action": [
                "dynamodb:List*",
                "dynamodb:DescribeReservedCapacity*",
                "dynamodb:DescribeLimits",
                "dynamodb:DescribeTimeToLive"
            ],
            "Resource": "*"
        },
        {
            "Sid": "SpecificTable",
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGet*",
                "dynamodb:DescribeStream",
                "dynamodb:DescribeTable",
                "dynamodb:Get*",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:BatchWrite*",
                "dynamodb:CreateTable",
                "dynamodb:Delete*",
                "dynamodb:Update*",
                "dynamodb:PutItem"
            ],
            "Resource": "arn:aws:dynamodb:*:*:table/userflow"
        }
    ]
}
```
