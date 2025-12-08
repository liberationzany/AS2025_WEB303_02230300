# Test Script for Student Cafe Microservices
# Run this after starting all services with docker-compose up

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Student Cafe Microservices Test Script" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8080"

# Function to make API calls
function Invoke-ApiCall {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Body = $null
    )
    
    try {
        if ($Body) {
            $response = Invoke-RestMethod -Method $Method -Uri "$baseUrl$Endpoint" -ContentType "application/json" -Body $Body
        } else {
            $response = Invoke-RestMethod -Method $Method -Uri "$baseUrl$Endpoint"
        }
        return $response
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
        return $null
    }
}

# Test 1: Create Users
Write-Host "Test 1: Creating users..." -ForegroundColor Yellow
$user1 = Invoke-ApiCall -Method POST -Endpoint "/api/users" -Body '{"name": "John Doe", "email": "john@example.com"}'
$user2 = Invoke-ApiCall -Method POST -Endpoint "/api/users" -Body '{"name": "Jane Smith", "email": "jane@example.com"}'

if ($user1) {
    Write-Host "✓ Created user: $($user1.name)" -ForegroundColor Green
    $userId = $user1.ID
} else {
    Write-Host "✗ Failed to create user" -ForegroundColor Red
    exit 1
}

if ($user2) {
    Write-Host "✓ Created user: $($user2.name)" -ForegroundColor Green
}

Write-Host ""

# Test 2: Create Menu Items
Write-Host "Test 2: Creating menu items..." -ForegroundColor Yellow
$menuItems = @(
    '{"name": "Coffee", "description": "Hot coffee", "price": 2.50}',
    '{"name": "Tea", "description": "Hot tea", "price": 1.50}',
    '{"name": "Sandwich", "description": "Fresh sandwich", "price": 5.00}',
    '{"name": "Salad", "description": "Green salad", "price": 4.50}'
)

$menuItemIds = @()
foreach ($item in $menuItems) {
    $created = Invoke-ApiCall -Method POST -Endpoint "/api/menu" -Body $item
    if ($created) {
        Write-Host "✓ Created menu item: $($created.name) - `$$($created.price)" -ForegroundColor Green
        $menuItemIds += $created.ID
    }
}

Write-Host ""

# Test 3: View Menu
Write-Host "Test 3: Retrieving menu..." -ForegroundColor Yellow
$menu = Invoke-ApiCall -Method GET -Endpoint "/api/menu"
if ($menu) {
    Write-Host "✓ Retrieved $($menu.Count) menu items" -ForegroundColor Green
    foreach ($item in $menu) {
        Write-Host "  - $($item.name): `$$($item.price)" -ForegroundColor Gray
    }
} else {
    Write-Host "✗ Failed to retrieve menu" -ForegroundColor Red
}

Write-Host ""

# Test 4: Get User
Write-Host "Test 4: Retrieving user..." -ForegroundColor Yellow
$retrievedUser = Invoke-ApiCall -Method GET -Endpoint "/api/users/$userId"
if ($retrievedUser) {
    Write-Host "✓ Retrieved user: $($retrievedUser.name) ($($retrievedUser.email))" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to retrieve user" -ForegroundColor Red
}

Write-Host ""

# Test 5: Create Order (Inter-service communication)
Write-Host "Test 5: Creating order (tests inter-service communication)..." -ForegroundColor Yellow
$orderBody = "{`"user_id`": $userId, `"items`": [{`"menu_item_id`": 1, `"quantity`": 2}, {`"menu_item_id`": 3, `"quantity`": 1}]}"
$order = Invoke-ApiCall -Method POST -Endpoint "/api/orders" -Body $orderBody

if ($order) {
    Write-Host "✓ Created order #$($order.ID) for user $($order.user_id)" -ForegroundColor Green
    Write-Host "  Status: $($order.status)" -ForegroundColor Gray
    Write-Host "  Items: $($order.order_items.Count)" -ForegroundColor Gray
    
    $total = 0
    foreach ($item in $order.order_items) {
        $subtotal = $item.quantity * $item.price
        $total += $subtotal
        Write-Host "    - Item #$($item.menu_item_id): $($item.quantity)x @ `$$($item.price) = `$$subtotal" -ForegroundColor Gray
    }
    Write-Host "  Total: `$$total" -ForegroundColor Cyan
} else {
    Write-Host "✗ Failed to create order" -ForegroundColor Red
}

Write-Host ""

# Test 6: View All Orders
Write-Host "Test 6: Retrieving all orders..." -ForegroundColor Yellow
$orders = Invoke-ApiCall -Method GET -Endpoint "/api/orders"
if ($orders) {
    Write-Host "✓ Retrieved $($orders.Count) order(s)" -ForegroundColor Green
    foreach ($ord in $orders) {
        Write-Host "  - Order #$($ord.ID): User $($ord.user_id), Status: $($ord.status), Items: $($ord.order_items.Count)" -ForegroundColor Gray
    }
} else {
    Write-Host "✗ Failed to retrieve orders" -ForegroundColor Red
}

Write-Host ""
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "All tests completed!" -ForegroundColor Green
Write-Host ""
Write-Host "Check Consul UI: http://localhost:8500" -ForegroundColor Cyan
Write-Host "View logs: docker-compose logs -f order-service" -ForegroundColor Cyan
Write-Host ""
