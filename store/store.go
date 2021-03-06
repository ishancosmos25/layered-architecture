package store

import "C"
import (
	"database/sql"
	"fmt"
	"layres/entities"
)

type CustomerStore struct {
	db *sql.DB
}


func (c CustomerStore)CloseDb(){
	c.db.Close()
}
func New() CustomerStore {
	var db, err = sql.Open("mysql", "root:Manish@123Sharma@/Customer_services")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return CustomerStore {db: db}
}

func (c CustomerStore) GetCustomerBYId(id int) (entities.Customer, error) {
	rows, err := c.db.Query("select * from customer inner join address on customer.id=address.cid and customer.id=? ", id)
	if err != nil {
		return entities.Customer{}, err
	}

	var cust entities.Customer

	for rows.Next() {
		rows.Scan(&cust.Id, &cust.Name, &cust.Dob, &cust.Add.Id, &cust.Add.StreetName, &cust.Add.City, &cust.Add.State, &cust.Add.CustomerId)
	}

	return cust, nil
}

func (c CustomerStore) GetCustomerByName(name string) (entities.Customer, error) {
	rows, err := c.db.Query("select * from customer inner join address on customer.id=address.cid where customer.name=? ", name)
	if err != nil {
		return entities.Customer{}, err
	}
	fmt.Println(name)
	var cust entities.Customer

	for rows.Next() {
		rows.Scan(&cust.Id, &cust.Name, &cust.Dob, &cust.Add.Id, &cust.Add.StreetName, &cust.Add.City, &cust.Add.State, &cust.Add.CustomerId)
	}
	fmt.Println(cust)
	return cust, nil
}

func (c CustomerStore) CreateCustomer(cust entities.Customer) (entities.Customer,error){
	var info[] interface{}
	query:=`insert into customer (name,dob) values(?,?)`
	if cust.Name=="" || cust.Dob==""{
		return entities.Customer{},nil
	}

	info=append(info,&cust.Name)
	info=append(info,&cust.Dob)

	row,_:=c.db.Exec(query,info...)
	query=`insert into address (street_name,city,state,cid) values(?,?,?,?)`
	var addr[] interface{}
	if cust.Add.StreetName=="" || cust.Add.City=="" || cust.Add.State==""{
		return entities.Customer{},nil
	}
	addr=append(addr,&cust.Add.StreetName)
	addr=append(addr,&cust.Add.City)
	addr=append(addr,&cust.Add.State)

	id,ok1:=row.LastInsertId()
	if ok1!=nil{
		return entities.Customer{},nil
	}
	addr=append(addr,id)
	_,ok:=c.db.Exec(query,addr...)
	if ok!=nil{
		panic(ok)
	}
	query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`


	newRow,_:=c.db.Query(query,id)
	var detail entities.Customer
	for newRow.Next() {
		newRow.Scan(&detail.Id, &detail.Name, &detail.Dob, &detail.Add.Id, &detail.Add.StreetName, &detail.Add.City, &detail.Add.State, &detail.Add.CustomerId)
	}
	return detail,nil
}

func (c CustomerStore) GetCustomer() ([]entities.Customer,error){
	query:=`select * from customer inner join address on customer.id=address.cid`
	rows,ok:=c.db.Query(query)
	if ok!=nil {
		panic(ok)
	}

	var response []entities.Customer

	defer rows.Close()

	for rows.Next() {
		var detail entities.Customer
		ok = rows.Scan(&detail.Id,&detail.Name,&detail.Dob,&detail.Add.Id,&detail.Add.StreetName,&detail.Add.City,&detail.Add.State,&detail.Add.CustomerId)
		response = append(response, detail)
	}
	return response,nil
}

func (c CustomerStore) RemoveCustomer(id int) error{
	var info[] interface{}
	info=append(info,id)

	query := `delete from customer where id=?`
	_, ok:= c.db.Exec(query, info...)
	if ok!=nil{
		return ok
	}
	return nil
}

func (c CustomerStore) UpdateCustomer (customer entities.Customer,id int) (entities.Customer,error){
	if customer.Name!=""{
		query:=`update customer set`
		var info [] interface{}
		query+=" name=?"
		info=append(info,customer.Name)
		query+=" where customer.id=?"
		info=append(info,id)
		_,er:=c.db.Exec(query,info...)

		if er!=nil{
			return entities.Customer{},er
		}
	}

	check:=entities.Address{}
	if  customer.Add!=check {
		query := `update address set `
		var idd []interface{}
		if customer.Add.StreetName != "" {
			idd = append(idd, customer.Add.StreetName)
			query += " street_name=?,"
		}

		if customer.Add.City != "" {
			idd = append(idd, customer.Add.City)
			query += " city=?,"
		}

		if customer.Add.State != "" {
			idd = append(idd, customer.Add.State)
			query += " state=?,"
		}
		query=query[:len(query)-1]
		query += " where address.cid=?"
		idd = append(idd, id)
		_, ok1 := c.db.Exec(query, idd...)

		if ok1 != nil {
			return entities.Customer{},ok1
		}
	}

	query:=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
	rows,_:=c.db.Query(query,id)
	var detail entities.Customer
	for rows.Next(){
		rows.Scan(&detail.Id,&detail.Name,&detail.Dob,&detail.Add.Id,&detail.Add.StreetName,&detail.Add.City,&detail.Add.State,&detail.Add.CustomerId)
	}
	if detail.Id==0{
		return entities.Customer{},nil
	}
	return detail,nil
}

