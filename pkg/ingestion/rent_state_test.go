package ingestion

import (
	"fmt"
	"testing"

	"github.com/dovydasdo/clerk/generated/ad"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRentStateUpdate(t *testing.T) {
	state := RentState{
		statsCity:      make(map[string]*cityStats, 0),
		totalProcessed: 0,
	}

	nAds := 20
	ads := make([]*ad.Ad, nAds)
	cityCount := 2

	prices := 0
	footage := 0.0

	prices_b := 0
	footage_b := 0.0

	for i := 0; i < nAds; i++ {
		var city string
		if i%2 == 0 {
			city = "city_a"
		} else {
			city = "city_b"
		}

		ad := ad.Ad{
			Id:             uint32(i),
			CreatedAt:      &timestamppb.Timestamp{},
			UpdatedAt:      &timestamppb.Timestamp{},
			DeletedAt:      &timestamppb.Timestamp{},
			City:           city,
			Date:           &timestamppb.Timestamp{},
			Stars:          0,
			Title:          "",
			Address:        "",
			Footage:        float64(i) * 2,
			Rooms:          0,
			Floor:          0,
			BuildingFloors: 0,
			Specifications: "",
			Price:          int32(i) * 10,
			Premium:        false,
			AdId:           "",
			AdIdUi:         0,
			Source:         "",
			Url:            "",
		}

		if city == "city_a" {
			prices += int(ad.Price)
			footage += ad.Footage
		}

		if city == "city_b" {
			prices_b += int(ad.Price)
			footage_b += ad.Footage
		}

		ads[i] = &ad
	}

	adsPerCity := nAds / 2
	avgPrice := prices / adsPerCity
	avgFootage := footage / float64(adsPerCity)
	avgPricePerSq := float64(avgPrice) / avgFootage

	avgPrice_b := prices_b / adsPerCity
	avgFootage_b := footage_b / float64(adsPerCity)
	avgPricePerSq_b := float64(avgPrice_b) / avgFootage_b

	for _, ad := range ads {
		UpdateRentState(ad, &state)
	}

	if len(state.statsCity) != cityCount {
		t.Error(fmt.Sprintf("expected %v cities, found: %v", cityCount, len(state.statsCity)))
	}

	if state.totalProcessed != nAds {
		t.Error(fmt.Sprintf("expected %v ads to process, found: %v", nAds, state.totalProcessed))
	}

	if state.statsCity["city_a"].avgPrice != avgPrice {
		t.Error(fmt.Sprintf("expected %v avg price, found: %v", avgPrice, state.statsCity["city_a"].avgPrice))
	}

	if state.statsCity["city_b"].avgPrice != avgPrice_b {
		t.Error(fmt.Sprintf("expected %v avg price, found: %v", avgPrice_b, state.statsCity["city_b"].avgPrice))
	}

	if state.statsCity["city_a"].avgFootage != avgFootage {
		t.Error(fmt.Sprintf("expected %v avg footage, found: %v", avgFootage, state.statsCity["city_a"].avgFootage))
	}

	if state.statsCity["city_b"].avgFootage != avgFootage_b {
		t.Error(fmt.Sprintf("expected %v avg footage, found: %v", avgFootage_b, state.statsCity["city_b"].avgFootage))
	}

	if state.statsCity["city_a"].avgPricePerSq != float64(avgPricePerSq) {
		t.Error(fmt.Sprintf("expected %v avg pps, found: %v", avgPricePerSq, state.statsCity["city_a"].avgPricePerSq))
	}

	if state.statsCity["city_b"].avgPricePerSq != float64(avgPricePerSq_b) {
		t.Error(fmt.Sprintf("expected %v avg pps, found: %v", avgPrice_b, state.statsCity["city_b"].avgPricePerSq))
	}
}
