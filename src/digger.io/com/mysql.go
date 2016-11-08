package com

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strings"
    "reflect"
    "errors"
)

type MysqlCom struct {
    db      *sql.DB
}

type MysqlConfig struct {
    Username    string
    Password    string
    Host        string
    Database    string
    Charset     string
    MaxOpen     int
    MaxIdle     int
}


func (s *MysqlCom) InitConn(config *MysqlConfig) *sql.DB {
    conStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", config.Username, config.Password, config.Host, config.Database, config.Charset)
    s.db, _ = sql.Open("mysql", conStr)
    s.db.SetMaxOpenConns(config.MaxOpen)
    s.db.SetMaxIdleConns(config.MaxIdle)
    s.db.Ping()
    return s.db
}

func (s *MysqlCom) GetConn() *sql.DB {
    return s.db
}

func (s *MysqlCom) Insert(table string, data map[string]string) int64 {
    // format datea
    field, value := s.sqlEscapeMap(data)
    stmt, err := s.GetConn().Prepare(fmt.Sprintf(`REPLACE INTO %s (%s) values (%s)`, table, field, value))
    if err != nil {
        panic(err.Error())
    }
    // exec sql
    res, err := stmt.Exec()
    if err != nil {
        panic(err.Error())
    }
    // get last id
    id, err := res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    return id
}

func (s *MysqlCom) Query(sql string, args []interface{}) ([]map[string]interface{}){
    // Execute the query
    rows, err := s.GetConn().Query(sql, args...)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Make a slice for the values
    values := make([]interface{}, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    var results []map[string]interface{}
    // Fetch rows
    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
        // Here we just print each column as a string.
        var row = make(map[string]interface{})
        for i, _ := range values {
            row[strings.Title(columns[i])] = values[i]
        }
        results = append(results, row)
    }
    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    return results
}

func (s *MysqlCom) QueryRow(sql string, args []interface{}) (map[string]interface{}){
    // Execute the query
    rows, err := s.GetConn().Query(sql, args...)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Make a slice for the values
    values := make([]interface{}, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    // Fetch rows
    rows.Next()
    // get RawBytes from data
    err = rows.Scan(scanArgs...)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    // Here we just print each column as a string.
    var row = make(map[string]interface{})
    for i, _ := range values {
        row[strings.Title(columns[i])] = values[i]
    }
    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    return row
}

func (s *MysqlCom) FillStruct(out interface{}, m map[string]interface{}) error{
    for k, v := range m {
        err := s.SetField(out, k, v)
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *MysqlCom) SetField(obj interface{}, name string, value interface{}) error {
    structValue := reflect.ValueOf(obj).Elem()
    structFieldValue := structValue.FieldByName(name)

    if !structFieldValue.IsValid() {
        return fmt.Errorf("No such field: %s in obj", name)
    }
    if !structFieldValue.CanSet() {
        return fmt.Errorf("Cannot set %s field value", name)
    }

    structFieldType := structFieldValue.Type()
    val := reflect.ValueOf(value)

    if structFieldType != val.Type() && structFieldType.String() == "string" && val.Type().String() == "[]uint8"{
        val = reflect.ValueOf(string(value.([]byte)))
    } else if structFieldType != val.Type() {
        fmt.Println(name, val.Type(), structFieldType)
        return errors.New("Provided value type didn't match obj field type")
    }

    structFieldValue.Set(val)
    return nil
}

func (s *MysqlCom) sqlEscapeMap(data map[string]string) (string, string) {
    field := ""
    value := ""
    for k, v := range data {
        field += fmt.Sprintf("`%s`,", k)
        value += fmt.Sprintf("\"%s\",", strings.Replace(v, "\"", "\\\"", -1))
    }
    return strings.Trim(field, ","), strings.Trim(value, ",")
}
