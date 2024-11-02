package actions

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Order struct{
  CustomerName string   `dynamodbav:"CustomerName" json:"CustomerName"`
  OrderID      int      `dynamodbav:"OrderID" json:"OrderID"`
  OrderStatus  string   `dynamodbav:"OrderStatus" json:"OrderStatus"`
  OrderDate    string   `dynamodbav:"OrderDate" json:"OrderDate"`
  TotalAmount  float32  `dynamodbav:"TotalAmount" json:"TotalAmount"`
  Items        []string `dynamodbav:"Items" json:"Items"`
}

func (order Order) GetKey() map[string]types.AttributeValue {
	customerName, err := attributevalue.Marshal(order.CustomerName)
	if err != nil {
		log.Fatal(err)
	}

	orderID, err := attributevalue.Marshal(order.OrderID)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{"CustomerName": customerName, "OrderID": orderID}
}

func (order Order) String() string {
	return fmt.Sprintf("\t%d\t%s\t%s\t%f\t%s\t%v\n", 
	      order.OrderID, order.CustomerName, order.OrderDate, order.TotalAmount, order.OrderStatus, order.Items)
}

