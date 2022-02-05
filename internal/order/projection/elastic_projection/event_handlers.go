package elastic_projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (o *elasticProjection) onOrderCreate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleOrderCreateEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderCreatedEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:      aggregate.GetOrderAggregateID(evt.AggregateID),
		ShopItems:    eventData.ShopItems,
		Created:      true,
		AccountEmail: eventData.AccountEmail,
		TotalPrice:   aggregate.GetShopItemsTotalPrice(eventData.ShopItems),
	}

	return o.elasticRepository.IndexOrder(ctx, op)
}

func (o *elasticProjection) onOrderPaid(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleOrderPaidEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var payment models.Payment
	if err := evt.GetJsonData(&payment); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Paid: true, Payment: payment}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) onSubmit(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleSubmitEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) onUpdate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleUpdateEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderUpdatedEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), ShopItems: eventData.ShopItems}
	op.TotalPrice = aggregate.GetShopItemsTotalPrice(eventData.ShopItems)
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) onCancel(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onCancel")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderCanceledEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:      aggregate.GetOrderAggregateID(evt.AggregateID),
		Canceled:     true,
		Delivered:    false,
		CancelReason: eventData.CancelReason,
	}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) onDelivered(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onDelivered")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderDeliveredEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:       aggregate.GetOrderAggregateID(evt.AggregateID),
		Delivered:     true,
		DeliveredTime: eventData.DeliveryTimestamp,
	}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) onOrderDeliveryAddressUpdated(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onOrderDeliveryAddressUpdated")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderChangeDeliveryAddress
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:         aggregate.GetOrderAggregateID(evt.AggregateID),
		DeliveryAddress: eventData.DeliveryAddress,
	}
	return o.elasticRepository.UpdateOrder(ctx, op)

}
