package csv_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Drelf2018/csv"
)

type QTime time.Time

func (q QTime) MarshalCSV() (string, error) {
	return (time.Time)(q).Format("2006/1/2 15:04"), nil
}

func (q *QTime) UnmarshalCSV(s string) error {
	t, err := time.Parse("2006/1/2 15:04", s)
	if err != nil {
		return err
	}
	*q = QTime(t)
	return nil
}

func (q QTime) String() string {
	return time.Time(q).Format("2006/01/02 15:04")
}

var _ csv.Marshaler = (*QTime)(nil)
var _ csv.Unmarshaler = (*QTime)(nil)

type Tags []string

func (t Tags) MarshalCSV() (string, error) {
	if len(t) == 0 {
		return "", nil
	} else if len(t) == 1 {
		return t[0], nil
	}
	return "#" + strings.Join(t, "#"), nil
}

func (t *Tags) UnmarshalCSV(s string) error {
	for _, v := range strings.Split(s, "#") {
		if v != "" {
			*t = append(*t, v)
		}
	}
	return nil
}

func (t Tags) String() string {
	return strings.Join(t, "/")
}

var _ csv.Marshaler = (*Tags)(nil)
var _ csv.Unmarshaler = (*Tags)(nil)

type Test struct {
	A1  *QTime `csv:"时间"`
	A2  string `csv:"分类"`
	A3  string `csv:"类型"`
	A4  []byte `csv:"金额"`
	A5  string `csv:"账户1"`
	A6  string `csv:"账户2"`
	A7  string `csv:"备注"`
	B1  int
	A8  string `csv:"账单标记"`
	A9  string `csv:"手续费"`
	A10 string `csv:"优惠券"`
	A11 *Tags  `csv:"标签"`
	A12 string `csv:"账单图片"`
}

func (Test) OrderedCSV() []string {
	return []string{"时间", "类型", "金额", "备注", "标签"}
}

var _ csv.Ordered = (*Test)(nil)

var data = []byte(`时间,分类,类型,金额,账户1,账户2,备注,账单标记,手续费,优惠券,标签,账单图片
2022/7/12 21:20,三餐,支出,28.58,微信,,去肯德基吃汉堡（此行数据是示例，可以删除）,不计收支,,,种草,http://billimg.qianjiapp.com/202006300908267611e24759b6c1c8781361af861!webporigin
2022/7/8 22:15,工资,收入,1000,,,7月份工资（此行数据是示例，可以删除）,,,,#老婆#旅行,
2022/7/8 10:10,,转账,200,支付宝,招商银行卡,支付宝提现1000元到银行卡（此行数据是示例，可以删除）,,,,#冲动消费#拔草,
2022/7/8 22:15,奶茶,支出,12.5,,,茶百道（此行数据是示例，可以删除）,不计收支&预算,,,,`)

func TestUnmarshal(t *testing.T) {
	records, err := csv.Unmarshal[Test](data)
	if err != nil {
		t.Fatal(err)
	}
	for i, record := range records {
		t.Logf("record#%d: %v\n", i, record)
	}
}

func TestMarshal(t *testing.T) {
	records, err := csv.Unmarshal[Test](data)
	if err != nil {
		t.Fatal(err)
	}
	// use csv.Ordered
	p, err := csv.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(p))
}
