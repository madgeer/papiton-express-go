# Event Flow

## 1. Prinsip
Tracking dan Notification adalah consumer utama untuk event eksternal. Namun, Order Service juga bertindak sebagai consumer untuk event dari layanan lain guna melakukan pembaruan status order secara asinkron.

## 2. Event Utama
| Event | Publisher | Consumer |
|------|-----------|----------|
| OrderCreated | Order Service | Payment Service, Tracking Service, Notification Service |
| PaymentVerified | Payment Service | Shipping Service, Order Service, Tracking Service, Notification Service |
| CourierAssigned | Shipping Service | Order Service, Tracking Service, Notification Service |
| PackagePickedUp | Shipping Service | Order Service, Tracking Service, Notification Service |
| WarehouseIn | Warehouse Service | Order Service, Tracking Service |
| WarehouseSorted | Warehouse Service | Tracking Service |
| WarehouseOut | Warehouse Service | Order Service, Tracking Service, Notification Service |
| DeliveryCompleted | Shipping Service | Order Service, Tracking Service, Notification Service |
| OrderCompleted | Order Service | Tracking Service, Notification Service |

## 3. Topic Kafka
- order.events (OrderCreated, OrderCompleted)
- payment.events (PaymentVerified)
- shipping.events (CourierAssigned, PackagePickedUp, DeliveryCompleted)
- warehouse.events (WarehouseIn, WarehouseSorted, WarehouseOut)

## 4. Alur Event
```mermaid
flowchart TD
    classDef service fill:#e1f5fe,stroke:#039be5,stroke-width:2px;
    classDef kafka fill:#fff3e0,stroke:#ffb74d,stroke-width:2px;

    %% Services
    OrderSvc[Order Service]:::service
    PaymentSvc[Payment Service]:::service
    ShippingSvc[Shipping Service]:::service
    WarehouseSvc[Warehouse Service]:::service
    TrackingSvc[Tracking Service]:::service
    NotifSvc[Notification Service]:::service

    %% Kafka Topics
    TopicOrder[Topic: order.events]:::kafka
    TopicPayment[Topic: payment.events]:::kafka
    TopicShipping[Topic: shipping.events]:::kafka
    TopicWarehouse[Topic: warehouse.events]:::kafka

    %% Event Publishing
    OrderSvc -->|Publish: OrderCreated, OrderCompleted| TopicOrder
    PaymentSvc -->|Publish: PaymentVerified| TopicPayment
    ShippingSvc -->|Publish: CourierAssigned, PackagePickedUp, DeliveryCompleted| TopicShipping
    WarehouseSvc -->|Publish: WarehouseIn, WarehouseSorted, WarehouseOut| TopicWarehouse

    %% Event Consuming
    TopicOrder --> PaymentSvc
    TopicOrder --> TrackingSvc
    TopicOrder --> NotifSvc

    TopicPayment --> ShippingSvc
    TopicPayment --> OrderSvc
    TopicPayment --> TrackingSvc
    TopicPayment --> NotifSvc

    TopicShipping --> OrderSvc
    TopicShipping --> TrackingSvc
    TopicShipping --> NotifSvc

    TopicWarehouse --> OrderSvc
    TopicWarehouse --> TrackingSvc
    TopicWarehouse --> NotifSvc
```

