package css

type AddressPart struct {
	element string
	ids     []string
	classes []string
}

type Address struct {
	AddressParts []AddressPart
}

type Style struct {
	Key   string
	Value string
}

type Stylesheet struct {
	Addresses []Address
	Styles    []Style
}
