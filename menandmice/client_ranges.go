package menandmice

type NextFreeAddressRespons struct {
	Result struct {
		Address string `json:"addres"`
	} `json:"result"`
}

func (c Mmclient) NextFreeAddress(addressRange, startaddress string, ping, excludeDHCP bool) (string, error) {
	var re NextFreeAddressRespons
	err := c.Get(&re, "Ranges/"+addressRange+"/NextFreeAddress", map[string]interface{}{
		"startaddress": startaddress,
		"ping":         ping,
		"excludeDHCP":  excludeDHCP,
	}, nil)
	return re.Result.Address, err
}
