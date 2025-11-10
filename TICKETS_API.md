# Tickets CRUD API Usage Examples

## Available Endpoints

### Create Ticket

```bash
POST http://localhost:4003/api/v1/tickets
Content-Type: application/json

{
  "ticket_id": "12345",
  "subject": "Test Issue",
  "description": "This is a test ticket",
  "status": 2,
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

### Get All Tickets (Paginated)

```bash
GET http://localhost:4003/api/v1/tickets?limit=10&offset=0
```

### Get Ticket by Database ID

```bash
GET http://localhost:4003/api/v1/tickets/1
```

### Get Ticket by Freshdesk Ticket ID

```bash
GET http://localhost:4003/api/v1/tickets/freshdesk/12345
```

### Get Tickets by Status

```bash
GET http://localhost:4003/api/v1/tickets/status/2?limit=10&offset=0
```

Status values (Freshdesk):

- 2 = Open
- 3 = Pending
- 4 = Resolved
- 5 = Closed

### Search Tickets

```bash
GET http://localhost:4003/api/v1/tickets/search?q=keyword&limit=10&offset=0
```

### Get Ticket Statistics

```bash
GET http://localhost:4003/api/v1/tickets/stats
```

### Update Ticket (Full)

```bash
PUT http://localhost:4003/api/v1/tickets/1
Content-Type: application/json

{
  "ticket_id": "12345",
  "subject": "Updated Subject",
  "description": "Updated description",
  "status": 3,
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T11:00:00Z"
}
```

### Update Ticket (Partial)

```bash
PATCH http://localhost:4003/api/v1/tickets/1
Content-Type: application/json

{
  "status": 4,
  "subject": "Only updating these fields"
}
```

### Delete Ticket (Soft Delete)

```bash
DELETE http://localhost:4003/api/v1/tickets/1
```

---

## Using the Repository Directly in Code

```go
package main

import (
    "github.com/DerylFeyza/freshdesk-automation/models"
    "github.com/DerylFeyza/freshdesk-automation/repository"
)

func main() {
    // Initialize repository
    repo := repository.NewTicketRepository()

    // Create a ticket
    ticket := &models.Tickets{
        Ticket_id:   "12345",
        Subject:     "New Issue",
        Description: "Issue description",
        Status:      2,
        Created_at:  "2025-11-10T10:00:00Z",
        Updated_at:  "2025-11-10T10:00:00Z",
    }
    err := repo.Create(ticket)

    // Find by ID
    ticket, err := repo.FindByID(1)

    // Find by Freshdesk ticket ID
    ticket, err := repo.FindByTicketID("12345")

    // Get all tickets with pagination
    tickets, err := repo.FindAll(10, 0) // limit=10, offset=0

    // Find by status
    openTickets, err := repo.FindByStatus(2, 10, 0)

    // Update
    ticket.Status = 3
    err = repo.Update(ticket)

    // Update specific fields
    updates := map[string]interface{}{
        "status":  4,
        "subject": "Updated",
    }
    err = repo.UpdateFields(1, updates)

    // Delete (soft delete)
    err = repo.Delete(1)

    // Hard delete (permanent)
    err = repo.HardDelete(1)

    // Check if exists
    exists, err := repo.ExistsByTicketID("12345")

    // Search
    results, err := repo.Search("keyword", 10, 0)

    // Count
    total, err := repo.Count()
    openCount, err := repo.CountByStatus(2)

    // Create or Update (upsert)
    err = repo.CreateOrUpdate(ticket)

    // Batch create
    tickets := []models.Tickets{ticket1, ticket2, ticket3}
    err = repo.BatchCreate(tickets)
}
```

## Response Examples

### Success Response (Single Ticket)

```json
{
	"ID": 1,
	"CreatedAt": "2025-11-10T10:00:00Z",
	"UpdatedAt": "2025-11-10T10:00:00Z",
	"DeletedAt": null,
	"ticket_uuid": "550e8400-e29b-41d4-a716-446655440000",
	"ticket_id": "12345",
	"subject": "Test Issue",
	"description": "This is a test ticket",
	"status": 2,
	"created_at": "2025-11-10T10:00:00Z",
	"updated_at": "2025-11-10T10:00:00Z"
}
```

### Success Response (Multiple Tickets)

```json
{
  "tickets": [...],
  "limit": 10,
  "offset": 0
}
```

### Error Response

```json
{
	"error": "ticket not found"
}
```
