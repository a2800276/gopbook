package gopbook

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// representation of the pass.json meta-data
type PassMetaData struct {
	Description        string `json:"description"`
	FormatVersion      uint   `json:"formatVersion"`
	OrganizationName   string `json:"organizationName"`
	PassTypeIdentifier string `json:"passTypeIdentifier"`
	SerialNumber       string `json:"serialNumber"`
	TeamIdentifier     string `json:"teamIdentifier"`
	// associated App Keys
	AppLaunchURL               string `json:"appLaunchURL"`
	AssociatedStoreIdentifiers string `json:"associatedStoreIdentifiers"`
	// companion App Keys
	UserInfo interface{} `json:"userInfo"`
	// Expiration Keys
	ExpirationDate W3Time `json:"expirationDate"`
	Voided         bool   `json:"voided"`
	// Relevance Keys
	Beacons      []BeaconDictionary   `json:"beacons"`
	Locations    []LocationDictionary `json:"locations"`
	MaxDistance  int                  `json:"maxDistance"`
	RelevantDate W3Time               `json:"relevantDate,omitempty"`
	// Style Keys
	BoardingPass PassStructureDictionary `json:"boardingPass"`
	Coupon       PassStructureDictionary `json:"coupon"`
	EventTicket  PassStructureDictionary `json:"eventTicket"`
	Generic      PassStructureDictionary `json:"generic"`
	StoreCard    PassStructureDictionary `json:"storeCard"`
	// Visual Appearance Keys
	Barcode            BarcodeDictionary `json:"barcode"`
	BackgroundColor    Color             `json:"backgroundColor"`
	ForegroundColor    Color             `json:"foregroundColor"`
	GroupingIdentifier string            `json:"groupingIdentifier"`
	LabelColor         Color             `json:"labelColor"`
	LogoText           string            `json:"logoText"`
	SuppressStripShine bool              `json:"suppressStripShine"`
	// Web Service Keys
	AuthenticationToken string `json:"authenticationToken"`
	WebServiceURL       string `json:"webServiceURL"`
}

type PassStructureDictionary struct {
	AuxiliaryFields []FieldDictionary `json:"auxiliaryFields"`
	BackFields      []FieldDictionary `json:"backFields"`
	HeaderFields    []FieldDictionary `json:"headerFields"`
	PrimaryFields   []FieldDictionary `json:"primaryFields"`
	SecondaryFields []FieldDictionary `json:"secondaryFields"`
	TransitType     TransitType       `json:"transitType"`
}

type BeaconDictionary struct {
	Major         uint16 `json:"major"`
	Minor         uint16 `json:"minor"`
	ProximityUUID string `json:"proximityUUID"`
	RelevantText  string `json:"relevantText"`
}

type LocationDictionary struct {
	Altitude     float64 `json:"altitude"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RelevantText float64 `json:"relevantText"`
}

type BarcodeDictionary struct {
	AltText         string `json:"altText"`
	Format          string `json:"format"`
	Message         string `json:"message"`
	MessageEncoding string `json:"messageEncoding"`
}

type FieldDictionary struct {
	AttributedValue   string             `json:"attributedValue"`
	ChangeMessage     string             `json:"changeMessage"`
	DataDetectorTypes []DataDetectorType `json:"dataDetectorTypes"`
	Key               string             `json:"key"`
	Label             string             `json:"label"`
	TextAlignment     TextAlignment      `json:"textAlignment"`
	Value             interface{}        `json:"value"`
}

type DateStyleField struct {
	DateStyle       DateStyle `json:"dateStyle"`
	IgnoresTimeZone bool      `json:"ignoresTimeZone"`
	IsRelative      bool      `json:"isRelative"`
	TimeStyle       string    `json:"timeStyle"`
}

type NumberStyleField struct {
	CurrencyCode string      `json:"currencyCode"`
	NumberStyle  NumberStyle `json:"numberStyle"`
}

type Color string

func MakeColor(r uint, g uint, b uint) Color {
	return Color(fmt.Sprintf("rgb(%d,%d,%d)", r, g, b))
}

type DataDetectorType string

const PhoneNumer DataDetectorType = "PKDataDectorTypePhoneNumer"
const Link DataDetectorType = "PKDataDectorTypeLink"
const Address DataDetectorType = "PKDataDectorTypeAddress"
const CalendarEvent DataDetectorType = "PKDataDectorTypeCalendarEvent"

type TextAlignment string

const Left TextAlignment = "PKTextAlignmentLeft"
const Center TextAlignment = "PKTextAlignmentCenter"
const Right TextAlignment = "PKTextAlignmentRight"
const Natural TextAlignment = "PKTextAlignmentNatural"

type DateStyle string

const DateStyleNone DateStyle = "NSDateFormatterNoStyle"
const DateStyleShort DateStyle = "NSDateFormatterShortStyle"
const DateStyleMedium DateStyle = "NSDateFormatterMediumStyle"
const DateStyleLong DateStyle = "NSDateFormatterLongStyle"
const DateStyleFull DateStyle = "NSDateFormatterFullStyle"

type NumberStyle string

const Decimal NumberStyle = "PKNumberStyleDecimal"
const Percent NumberStyle = "PKNumberStylePercent"
const Scientific NumberStyle = "PKNumberStyleScientifc"
const SpellOut NumberStyle = "PKNumberStyleSpellout"

type TransitType string

const TransitAir TransitType = "PKTransitTypeAir"
const TransitBoat TransitType = "PKTransitTypeBoat"
const TransitBus TransitType = "PKTransitTypeBus"
const TransitGeneric TransitType = "PKTransitTypeGeneric"
const TransitTrain TransitType = "PKTransitTypeTrain"

func LoadPassMetaData(reader io.Reader) (pb PassMetaData, err error) {
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&pb)
	return
}
func (p *PassMetaData) SavePassMeta(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(*p)
	return err
}

type W3Time time.Time

func (t *W3Time) UnmarshalJSON(data []byte) (err error) {
	format := "\"2006-01-02T15:04-07:00\""
	if ti, err := time.Parse(format, string(data)); err != nil {
		return err
	} else {
		t = (*W3Time)(&ti)
	}

	return
}
