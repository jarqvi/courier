package db

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jarqvi/courier/internal/dns"
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	Name string `gorm:"unique;not null;index"`

	ARecord  []byte `gorm:"type:json"`
	MXRecord []byte `gorm:"type:json"`

	RecordsSet    bool      `gorm:"default:false"`
	LastCheckedAt time.Time `gorm:"default:null"`

	Addresses []Address `gorm:"foreignKey:DomainID"`
	Users     []User    `gorm:"foreignKey:DomainID"`
}

func (d *Database) CreateDomain(name string) (*Domain, error) {
	ipv4, err := getHostIPV4()
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	domain := &Domain{
		Name: name,
		ARecord: []byte(fmt.Sprintf(`{
			"Name": "%s",
			"IPV4": "%s"
		}`, name, ipv4)),
		MXRecord: []byte(fmt.Sprintf(`{
			"Name": "%s",
			"Priority": 10,
			"MailServer": "%s"
		}`, name, name)),
	}

	if err := d.Create(domain).Error; err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return domain, nil
}

func getHostIPV4() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("failed to get host ipv4: %w", err)
	}

	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read host ipv4: %w", err)
	}

	return string(ip), nil
}

func (d *Database) GetDomain(name string) (*Domain, error) {
	var domain Domain
	if err := d.Where("name = ?", name).First(&domain).Error; err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return &domain, nil
}

func (d *Database) GetAllDomains() ([]Domain, error) {
	var domains []Domain
	if err := d.Find(&domains).Error; err != nil {
		return nil, fmt.Errorf("failed to get all domains: %w", err)
	}

	return domains, nil
}

func (d *Database) DeleteDomain(name string) error {
	domain, err := d.GetDomain(name)
	if err != nil {
		return fmt.Errorf("domain not found: %w", err)
	}

	if err := d.Delete(domain).Error; err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	return nil
}

type CheckDomainRecordErrors struct {
	ARecordError  error
	MXRecordError error

	CustomError error
}

func (c CheckDomainRecordErrors) HasError() bool {
	return c.ARecordError != nil || c.MXRecordError != nil || c.CustomError != nil
}

func (d *Database) CheckDomainRecords(name string, dns *dns.DNS) CheckDomainRecordErrors {
	domain, err := d.GetDomain(name)
	if err != nil {
		return CheckDomainRecordErrors{
			CustomError: fmt.Errorf("failed to check domain records: %w", err),
		}
	}

	return CheckDomainRecordErrors{
		ARecordError:  checkARecord(domain.ARecord, dns),
		MXRecordError: checkMXRecord(domain.MXRecord, dns),
	}
}

func checkARecord(record []byte, dns *dns.DNS) error {
	aRecord := make(map[string]any)
	json.Unmarshal(record, &aRecord)

	ips, err := dns.ResolveARecord(aRecord["Name"].(string))
	if err != nil {
		return fmt.Errorf("failed to resolve a record for %s: %w", aRecord["Name"].(string), err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("no a record found for %s", aRecord["Name"].(string))
	}

	for _, ip := range ips {
		if ip == aRecord["IPV4"].(string) {
			return nil
		}
	}

	return fmt.Errorf("a record %s does not match %s", ips, aRecord["IPV4"].(string))
}

func checkMXRecord(record []byte, dns *dns.DNS) error {
	mxRecord := make(map[string]any)
	json.Unmarshal(record, &mxRecord)

	mxs, err := dns.ResolveMXRecord(mxRecord["Name"].(string))
	if err != nil {
		return fmt.Errorf("failed to resolve mx record for %s: %w", mxRecord["Name"].(string), err)
	}

	if len(mxs) == 0 {
		return fmt.Errorf("no mx record found for %s", mxRecord["Name"].(string))
	}

	for _, mx := range mxs {
		if mx == mxRecord["MailServer"].(string) {
			return nil
		}
	}

	return fmt.Errorf("mx record %s does not match %s", mxs, mxRecord["MailServer"].(string))
}
