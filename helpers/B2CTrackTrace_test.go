package helpers

import (
	"testing"
)

func TestGetCourier(t *testing.T) {
	counter := 0

	if courier, _ := GetCourier("https://www.trackyourparcel.eu/?orderref=P2023-016189-dNYTb&postcode=9041&country=BEL"); courier == "be-post" {
		counter++
	}

	if courier, _ := GetCourier("https://www.trackyourparcel.eu/?orderref=P2023-016416-7ZQnl&postcode=6713+VC&country=NLD"); courier == "dhl" {
		counter++
	}

	if courier, _ := GetCourier("https://www.trackyourparcel.eu/?orderref=P2023-016411-7AdVz&postcode=12650&country=SWE"); courier == "dhl" {
		counter++
	}

	if courier, _ := GetCourier("https://www.trackyourparcel.eu/?orderref=P2023-016408-wLtqn&postcode=74321&country=DEU"); courier == "dhl" {
		counter++
	}

	if courier, _ := GetCourier("https://www.trackyourparcel.eu/?orderref=P2023-016361-Z43Te&postcode=170+00&country=CZE"); courier == "dpd" {
		counter++
	}

	//https://www.trackyourparcel.eu/?orderref=P2023-016343-BiqPt&postcode=82655&country=FIN

	if counter == 5 {
		return
	}

	t.Fatal("")
}
