package actions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TableBasics struct {
	DynamoDBClient *dynamodb.Client
	TableName      string
}

func(b TableBasics) AddOrder(order *Order) error {
	item, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}
	
	_, err = b.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(b.TableName),
		Item: item,
	})
	if err != nil {
		return err
	}

	return nil
}

func(b TableBasics) GetOrder(customerName string, orderID int) (Order, error) {
	order := Order{CustomerName: customerName, OrderID: orderID}
	res, err := b.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: order.GetKey(),
		TableName: aws.String(b.TableName),
	})
	if err != nil {
		return order, err
	}

	if res.Item == nil {
		return order, fmt.Errorf("order with customerName %s and orderID %d not found", customerName, orderID)
	}

	err = attributevalue.UnmarshalMap(res.Item, &order)
	if err != nil {
		return order, err
	}

	return order, nil
}

func(b TableBasics) UpdateOrder(order *Order) (map[string]string, error) {
	var attributeMap map[string]string

	update := expression.Set(expression.Name("OrderStatus"), expression.Value(order.OrderStatus))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return attributeMap, err
	}

	res, err := b.DynamoDBClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: order.GetKey(),
		TableName: aws.String(b.TableName),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression: expr.Update(),
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return attributeMap, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &attributeMap)
	if err != nil {
		return attributeMap, err
	}

	return attributeMap, err
}

func(b TableBasics) Query(customerName string) ([]Order, error) {
	var orders []Order

	keyEx := expression.Key("CustomerName").Equal(expression.Value(customerName))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return orders, err
	}

	res, err := b.DynamoDBClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String(b.TableName),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression: expr.KeyCondition(),
	})
	if err != nil {
		return orders, err
	}

	err = attributevalue.UnmarshalListOfMaps(res.Items, &orders)
	if err != nil {
		return orders, err
	}

	return orders, err
}

func(b TableBasics) ListOrdersByStatus(orderStatus string) ([]Order, error) {
	var orders []Order

	keyEx := expression.Key("OrderStatus").Equal(expression.Value(orderStatus))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return orders, err
	}

	res, err := b.DynamoDBClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String(b.TableName),
		IndexName: aws.String("OrderStatus-OrderDate-index"),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression: expr.KeyCondition(),
	})
	if err != nil {
		return orders, err
	}

	err = attributevalue.UnmarshalListOfMaps(res.Items, &orders)
	if err != nil {
		return orders, err
	}

	return orders, err
}