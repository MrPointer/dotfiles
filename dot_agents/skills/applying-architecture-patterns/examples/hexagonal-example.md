# Hexagonal Architecture Implementation Example

Complete Python implementation demonstrating ports and adapters pattern.

## Core Domain Service

```python
class OrderService:
    """Domain service - no infrastructure dependencies."""

    def __init__(
        self,
        order_repository: OrderRepositoryPort,
        payment_gateway: PaymentGatewayPort,
        notification_service: NotificationPort
    ):
        self.orders = order_repository
        self.payments = payment_gateway
        self.notifications = notification_service

    async def place_order(self, order: Order) -> OrderResult:
        # Business logic
        if not order.is_valid():
            return OrderResult(success=False, error="Invalid order")

        # Use ports (interfaces)
        payment = await self.payments.charge(
            amount=order.total,
            customer=order.customer_id
        )

        if not payment.success:
            return OrderResult(success=False, error="Payment failed")

        order.mark_as_paid()
        saved_order = await self.orders.save(order)

        await self.notifications.send(
            to=order.customer_email,
            subject="Order confirmed",
            body=f"Order {order.id} confirmed"
        )

        return OrderResult(success=True, order=saved_order)
```

## Ports (Interfaces)

```python
from abc import ABC, abstractmethod

class OrderRepositoryPort(ABC):
    @abstractmethod
    async def save(self, order: Order) -> Order:
        pass

class PaymentGatewayPort(ABC):
    @abstractmethod
    async def charge(self, amount: Money, customer: str) -> PaymentResult:
        pass

class NotificationPort(ABC):
    @abstractmethod
    async def send(self, to: str, subject: str, body: str):
        pass
```

## Adapters (Implementations)

### Production Adapter

```python
class StripePaymentAdapter(PaymentGatewayPort):
    """Primary adapter: connects to Stripe API."""

    def __init__(self, api_key: str):
        self.stripe = stripe
        self.stripe.api_key = api_key

    async def charge(self, amount: Money, customer: str) -> PaymentResult:
        try:
            charge = self.stripe.Charge.create(
                amount=amount.cents,
                currency=amount.currency,
                customer=customer
            )
            return PaymentResult(success=True, transaction_id=charge.id)
        except stripe.error.CardError as e:
            return PaymentResult(success=False, error=str(e))
```

### Test Adapter

```python
class MockPaymentAdapter(PaymentGatewayPort):
    """Test adapter: no external dependencies."""

    async def charge(self, amount: Money, customer: str) -> PaymentResult:
        return PaymentResult(success=True, transaction_id="mock-123")
```

## Key Takeaways

- **Domain Core**: Business logic with no infrastructure imports
- **Ports**: Abstract interfaces defining what the domain needs
- **Adapters**: Concrete implementations (Stripe, Mock, PostgreSQL, etc.)
- **Swappability**: Easily swap adapters for testing or different environments
