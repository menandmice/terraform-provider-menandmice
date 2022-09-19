package menandmice

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func tryGetString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)

	}
	return ""
}

func testAccCheckResourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ref set")
		}

		return nil
	}
}

// convert structure to map ignore `json:"omitempty"`
func toMap(item interface{}) (map[string]interface{}, error) {

	var properties map[string]interface{}
	serialized, err := json.Marshal(item)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serialized, &properties)

	if err != nil {
		return nil, err
	}
	return properties, nil
}

func ipv6AddressDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldIP := net.ParseIP(old)
	newIP := net.ParseIP(new)

	return oldIP.Equal(newIP)
}

// TODO rewrite think can be simpler
// code inspired by: https://github.com/apparentlymart/go-cidr/blob/v1.1.0/cidr/cidr.go
// AddressRange returns the first and last addresses in the given CIDR range.
func AddressRange(network *net.IPNet) (net.IP, net.IP) {
	// the first IP is easy
	firstIP := network.IP

	// the last IP is the network address OR NOT the mask address
	prefixLen, bits := network.Mask.Size()
	if prefixLen == bits {
		// Easy!
		// But make sure that our two slices are distinct, since they
		// would be in all other cases.
		lastIP := make([]byte, len(firstIP))
		copy(lastIP, firstIP)
		return firstIP, lastIP
	}

	firstIPInt, bits := ipToInt(firstIP)
	hostLen := uint(bits) - uint(prefixLen)
	lastIPInt := big.NewInt(1)
	lastIPInt.Lsh(lastIPInt, hostLen)
	lastIPInt.Sub(lastIPInt, big.NewInt(1))
	lastIPInt.Or(lastIPInt, firstIPInt)

	return firstIP, intToIP(lastIPInt, bits)
}

// code inspired by: https://github.com/apparentlymart/go-cidr/blob/v1.1.0/cidr/cidr.go
func ipToInt(ip net.IP) (*big.Int, int) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	if len(ip) == net.IPv4len {
		return val, 32
	} else if len(ip) == net.IPv6len {
		return val, 128
	} else {
		panic(fmt.Errorf("Unsupported address length %d", len(ip)))
	}
}

// code inspired by: https://github.com/apparentlymart/go-cidr/blob/v1.1.0/cidr/cidr.go
func intToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	// Pack our IP bytes into the end of the return array,
	// since big.Int.Bytes() removes front zero padding.
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return net.IP(ret)
}

func toError(diags diag.Diagnostics) error {

	if diags.HasError() {

		var message string
		for _, diag := range diags {
			message += diag.Summary + "\n"

		}
		return fmt.Errorf(message)

	}
	return nil
}
