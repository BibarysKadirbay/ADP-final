package unit

import (
	"encoding/json"
	"testing"

	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrderConfirmedEventJSONContract(t *testing.T) {
	event := usecase.OrderConfirmedEvent{
		OrderID: uuid.New(), RestaurantID: uuid.New(), CustomerID: uuid.New(),
		PickupAddress: "10 Abay Ave", DeliveryAddress: "20 Dostyk Ave", RouteDistanceKM: 4.5,
	}
	raw, err := json.Marshal(event)
	require.NoError(t, err)
	var decoded usecase.OrderConfirmedEvent
	require.NoError(t, json.Unmarshal(raw, &decoded))
	require.Equal(t, event.OrderID, decoded.OrderID)
	require.Equal(t, event.RouteDistanceKM, decoded.RouteDistanceKM)
}
