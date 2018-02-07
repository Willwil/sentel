//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package loader

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/whaler/pipeline/data"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlLoader struct {
	target    registry.RuleDataTarget
	productId string
	db        *sql.DB
}

func newMysqlLoader(c config.Config) (Loader, error) {
	target := c.MustValue("datatarget").(registry.RuleDataTarget)
	productId := c.MustString("productId")
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", target.Username, target.Password, target.DatabaseHost, target.DatabaseName)
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &mysqlLoader{
		db:        db,
		target:    target,
		productId: productId,
	}, nil
}

func (p *mysqlLoader) Name() string { return "mysql" }
func (p *mysqlLoader) Close()       { p.db.Close() }

func (p *mysqlLoader) Load(f *data.DataFrame) error {
	dd := f.PrettyJson()
	// prepar sql statement
	vals := []interface{}{}
	cols := []string{}
	s := fmt.Sprintf("INSERT INTO %s(", p.productId)
	for k, v := range dd {
		s = fmt.Sprintf("%s%s,", s, k)
		cols = append(cols, "?")
		vals = append(vals, v)
	}
	s = strings.TrimRight(s, ",")
	s += ")"
	s += "VALUES("
	for i := 0; i < len(cols); i++ {
		s += "?,"
	}
	s = strings.TrimRight(s, ",")
	s += ")"
	stmt, err := p.db.Prepare(s)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals)
	return err
}
