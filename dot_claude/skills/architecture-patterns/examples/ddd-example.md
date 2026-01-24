# Domain-Driven Design Implementation Example

Complete Python implementation demonstrating DDD tactical patterns.

## Value Objects (Immutable)

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class Email:
    """Value object: validated email."""
    value: str

    def __post_init__(self):
        if "@" not in self.value:
            raise ValueError("Invalid email")

@dataclass(frozen=True)
class Money:
    """Value object: amount with currency."""
    amount: int  # cents
    currency: str

    def add(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise ValueError("Currency mismatch")
        return Money(self.amount + other.amount, self.currency)
```

## Entities (With Identity)

```python
from typing import List

class Order:
    """Entity: has identity, mutable state."""

    def __init__(self, id: str, customer: Customer):
        self.id = id
        self.customer = customer
        self.items: List[OrderItem] = []
        self.status = OrderStatus.PENDING
        self._events: List[DomainEvent] = []

    def add_item(self, product: Product, quantity: int):
        """Business logic in entity."""
        item = OrderItem(product, quantity)
        self.items.append(item)
        self._events.append(ItemAddedEvent(self.id, item))

    def total(self) -> Money:
        """Calculated property."""
        return sum(item.subtotal() for item in self.items)

    def submit(self):
        """State transition with business rules."""
        if not self.items:
            raise ValueError("Cannot submit empty order")
        if self.status != OrderStatus.PENDING:
            raise ValueError("Order already submitted")

        self.status = OrderStatus.SUBMITTED
        self._events.append(OrderSubmittedEvent(self.id))
```

## Aggregates (Consistency Boundary)

```python
from typing import Optional

class Customer:
    """Aggregate root: controls access to entities."""

    def __init__(self, id: str, email: Email):
        self.id = id
        self.email = email
        self._addresses: List[Address] = []
        self._orders: List[str] = []  # Order IDs, not full objects

    def add_address(self, address: Address):
        """Aggregate enforces invariants."""
        if len(self._addresses) >= 5:
            raise ValueError("Maximum 5 addresses allowed")
        self._addresses.append(address)

    @property
    def primary_address(self) -> Optional[Address]:
        return next((a for a in self._addresses if a.is_primary), None)
```

## Domain Events

```python
from dataclasses import dataclass, field
from datetime import datetime

@dataclass
class OrderSubmittedEvent:
    order_id: str
    occurred_at: datetime = field(default_factory=datetime.now)
```

## Repository (Aggregate Persistence)

```python
class OrderRepository:
    """Repository: persist/retrieve aggregates."""

    async def find_by_id(self, order_id: str) -> Optional[Order]:
        """Reconstitute aggregate from storage."""
        pass

    async def save(self, order: Order):
        """Persist aggregate and publish events."""
        await self._persist(order)
        await self._publish_events(order._events)
        order._events.clear()
```

## Key Takeaways

- **Value Objects**: Immutable, equality by attributes, self-validating
- **Entities**: Identity-based equality, encapsulate business rules
- **Aggregates**: Consistency boundaries, accessed via aggregate root
- **Domain Events**: Record what happened for eventual consistency
- **Repositories**: Abstract persistence, work with aggregates
