package learnthephrase

import (
	"math/rand"
	"time"
)

type Pharse struct {
	pharses []string
}

func NewPharse() *Pharse {
	rand.Seed(time.Now().UnixNano())
	return &Pharse{
		pharses: []string{
			"— Скажи-ка, дядя, ведь недаром",
			"Москва, спаленная пожаром,",
			"Французу отдана?",
			"— Вот те на!",
			"— Ведь были ж схватки боевые,", "Да, говорят, еще какие!",
			"Недаром помнит вся Россия",
			"Про день Бородина!",
		},
	}
}
func (p Pharse) DoWork(payload map[string]string) (string, error) {
	return p.pharses[rand.Intn(len(p.pharses))], nil
}
