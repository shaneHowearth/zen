package repo

import (
	"fmt"
	"strconv"
	"strings"
)

var _ item = (*Organisation)(nil)

// Organisation -
type Organisation struct {
	ID            int      `json:"_id"`
	URL           string   `json:"url"`
	ExternalID    string   `json:"external_id"`
	Name          string   `json:"name"`
	DomainNames   []string `json:"domain_names"`
	CreatedAt     string   `json:"created_at"`
	Details       string   `json:"details"`
	SharedTickets bool     `json:"shared_tickets"`
	Tags          []string `json:"tags"`
}

// ToDTO - TODO
func (o *Organisation) ToDTO() map[string][]string {
	m := map[string][]string{}
	m["_id"] = []string{fmt.Sprintf("%d", o.ID)}
	m["url"] = []string{o.URL}
	m["external_id"] = []string{o.ExternalID}
	m["name"] = []string{o.Name}
	m["domain_names"] = o.DomainNames
	m["created_at"] = []string{o.CreatedAt}
	m["details"] = []string{o.Details}
	m["shared_tickets"] = []string{strconv.FormatBool(o.SharedTickets)}
	m["tags"] = o.Tags
	return m
}

// CreateIndex -
// Ignore "returns unexported type" linter complaint
// nolint:golint
func (o *Organisation) CreateIndex(in interface{}, name string) map[string]map[string][]item {
	d := in.([]*Organisation)
	// map[fieldname]map[fieldvalue][]*Organisation
	m := make(map[string]map[string][]item)
	m["_id"] = make(map[string][]item)
	m["url"] = make(map[string][]item)
	m["external_id"] = make(map[string][]item)
	m["name"] = make(map[string][]item)
	m["domain_names"] = make(map[string][]item)
	m["created_at"] = make(map[string][]item)
	m["details"] = make(map[string][]item)
	m["shared_tickets"] = make(map[string][]item)
	m["tags"] = make(map[string][]item)

	for i := range d {
		m["_id"][fmt.Sprintf("%d", d[i].ID)] = append(m["_id"][fmt.Sprintf("%d", d[i].ID)], d[i])
		m["url"][strings.ToLower(d[i].URL)] = append(m["url"][strings.ToLower(d[i].URL)], d[i])
		m["external_id"][strings.ToLower(d[i].ExternalID)] = append(m["external_id"][strings.ToLower(d[i].ExternalID)], d[i])
		m["name"][strings.ToLower(d[i].Name)] = append(m["name"][strings.ToLower(d[i].Name)], d[i])
		for _, domainName := range d[i].DomainNames {
			m["domain_names"][strings.ToLower(domainName)] = append(m["domain_names"][strings.ToLower(domainName)], d[i])
		}
		m["created_at"][strings.ToLower(d[i].CreatedAt)] = append(m["created_at"][strings.ToLower(d[i].CreatedAt)], d[i])
		m["details"][strings.ToLower(d[i].Details)] = append(m["details"][strings.ToLower(d[i].Details)], d[i])
		m["shared_tickets"][strconv.FormatBool(d[i].SharedTickets)] = append(m["shared_tickets"][strconv.FormatBool(d[i].SharedTickets)], d[i])
		for _, tag := range d[i].Tags {
			m["tags"][strings.ToLower(tag)] = append(m["tags"][strings.ToLower(tag)], d[i])
		}
	}
	return m
}
