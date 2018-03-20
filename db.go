package main

import (
	"database/sql"
	"fmt"
)

//Gets from DB all domains
func getDomains(db *sql.DB) ([]Domain, error) {
	statement := fmt.Sprintf("SELECT domains.id,domains.name,domains.expiration,contacts.first_name,contacts.last_name,contacts.email FROM domains inner join contacts on contacts.id=domains.owner_contact_id")
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //when all have finished
	domains := []Domain{}
	for rows.Next() {
		dom := Domain{Owner: new(Contact)}                                                                                                  //init pointer!!!
		if err := rows.Scan(&dom.ID, &dom.Name, &dom.Expiration, &dom.Owner.Firstname, &dom.Owner.Lastname, &dom.Owner.Email); err != nil { //short statement
			return nil, err
		}
		domains = append(domains, dom)
	}
	return domains, nil
}

//Gets from DB domain by ID
func (dom *Domain) getDomain(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT domains.name,domains.expiration,contacts.first_name,contacts.last_name,contacts.email FROM domains inner join contacts on contacts.id=domains.owner_contact_id WHERE domains.id=%d", dom.ID)
	return db.QueryRow(statement).Scan(&dom.Name, &dom.Expiration, &dom.Owner.Firstname, &dom.Owner.Lastname, &dom.Owner.Email)
}

//Create domain and connected contact in DB
func (dom *Domain) createDomain(db *sql.DB) error {
	//insert contacts row
	statement := fmt.Sprintf("INSERT INTO contacts (first_name, last_name, email) VALUES('%s', '%s', '%s')", dom.Owner.Firstname, dom.Owner.Lastname, dom.Owner.Email)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	//take the new contact id
	var ownerContactID int
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&ownerContactID)
	if err != nil {
		return err
	}
	//insert domain row
	statement = fmt.Sprintf("INSERT INTO domains (name, expiration,owner_contact_id) VALUES('%s', '%s', %d)", dom.Name, dom.Expiration, ownerContactID)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	//take the new domain id
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&dom.ID)
	if err != nil {
		return err
	}
	return nil
}

//update domain and connected contact in DB
func (dom *Domain) updateDomain(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE domains,contacts SET domains.name='%s', domains.expiration='%s',contacts.first_name='%s',contacts.last_name='%s',contacts.email='%s' WHERE contacts.id=domains.owner_contact_id and domains.id=%d", dom.Name, dom.Expiration, dom.Owner.Firstname, dom.Owner.Lastname, dom.Owner.Email, dom.ID)
	_, err := db.Exec(statement)
	return err
}

//delete domain and connected contact in DB
func (dom *Domain) deleteDomain(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE domains.*,contacts.* from domains,contacts WHERE contacts.id=domains.owner_contact_id and domains.id=%d", dom.ID)
	_, err := db.Exec(statement)
	return err
}
