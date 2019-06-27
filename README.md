# userflow

Lazypic 사용자 데이터 입니다.

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
$ userflow -rm -email woong@lazypic.org
```

### AWS DB권한 설정
DB접근 권한이 있는 계정에 아래 권한을 부여한다.

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