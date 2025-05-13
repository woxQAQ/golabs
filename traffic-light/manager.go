package trafficlight

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type car struct {
	id       int
	callback func(carID int, status lightStatus, m *manager)
}

func (c *car) run(id int) {
	fmt.Printf("Car %d: Running (Green Light), passing intersection.\n", id)
}

type manager struct {
	currentStatus  atomic.Value
	redDuration    time.Duration
	greenDuration  time.Duration
	yellowDuration time.Duration
	mu             sync.Mutex
	// Use a slice to maintain order of arrival
	// This slice will act as a queue for cars waiting at the red light
	// or cars that have arrived and are waiting for their first status update.
	registeredCars []*car
	nextCarID      int32
}

func newManager() *manager {
	m := &manager{
		redDuration:    10 * time.Second,
		greenDuration:  5 * time.Second, // Let's make green light a bit longer to see cars pass
		yellowDuration: 2 * time.Second,
		registeredCars: make([]*car, 0),
	}
	m.currentStatus.Store(LIGHT_STATUS_RED)
	return m
}

func (m *manager) RegisterCar(c *car) {
	m.mu.Lock()
	newID := int(atomic.AddInt32(&m.nextCarID, 1))
	c.id = newID
	m.registeredCars = append(m.registeredCars, c) // Add to the end of the queue
	initialStatus := m.getStatus()
	m.mu.Unlock() // Unlock before calling callback to avoid deadlock if callback calls manager

	fmt.Printf("Car %d arrived and registered. Current light: %v\n", c.id, initialStatus)
	// Notify the new car of the current status immediately
	go c.callback(c.id, initialStatus, m) // Pass manager so car can deregister itself if needed
}

// DeregisterCar removes a car by its ID from the registeredCars slice.
func (m *manager) DeregisterCar(carID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	for i, c := range m.registeredCars {
		if c.id == carID {
			// Remove element by slicing (preserves order of other elements)
			m.registeredCars = append(m.registeredCars[:i], m.registeredCars[i+1:]...)
			fmt.Printf("Car %d deregistered.\n", carID)
			found = true
			break
		}
	}
	if !found {
		// This might happen if a car is already deregistered by another concurrent callback
		// Or if the car was never registered with that ID, though less likely with atomic nextCarID
		// fmt.Printf("Car %d not found for deregistration (อาจจะถูกยกเลิกการลงทะเบียนไปแล้ว).\n", carID)
	}
}

func (m *manager) cycle(ctx context.Context) {
	ticker := time.NewTicker(m.redDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := m.getStatus()
			newStatus := current
			switch current {
			case LIGHT_STATUS_RED:
				newStatus = LIGHT_STATUS_GREEN
				ticker.Reset(m.greenDuration)
			case LIGHT_STATUS_GREEN:
				newStatus = LIGHT_STATUS_YELLOW
				ticker.Reset(m.yellowDuration)
			case LIGHT_STATUS_YELLOW:
				newStatus = LIGHT_STATUS_RED
				ticker.Reset(m.redDuration)
			}
			if newStatus != current {
				m.currentStatus.Store(newStatus)
				m.notifyCars(newStatus) // Notify all cars about the state change
			}
		}
	}
}

// notifyCars informs all registered cars about the light status change.
// If the light is green, it implies cars that can go, will go (and deregister in their callback).
func (m *manager) notifyCars(status lightStatus) {
	m.mu.Lock()
	// Create a copy of the current car list to iterate over.
	// This is important because car callbacks might modify the m.registeredCars slice (deregister).
	carsSnapshot := make([]*car, len(m.registeredCars))
	copy(carsSnapshot, m.registeredCars)
	m.mu.Unlock()

	fmt.Printf("Light changed to %v. Notifying %d cars.\n", status, len(carsSnapshot))

	for _, c := range carsSnapshot {
		// Check if the car still exists in the original list before notifying.
		// This is a safeguard if a car was deregistered between snapshot and this point by another goroutine,
		// though less likely with current lock scopes. More robust would be to check existence inside callback or before.
		// For simplicity here, we assume car.callback will handle if car.id is already deregistered.
		go c.callback(c.id, status, m)
	}
}

func (m *manager) getStatus() lightStatus {
	return m.currentStatus.Load().(lightStatus)
}

func Run(ctx context.Context) {
	m := newManager()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	// Car generation goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			// Generate cars periodically
			case <-time.After(time.Duration(r.Intn(2)+1) * time.Second): // Random interval 1-2 seconds
				newCar := &car{}
				newCar.callback = func(id int, status lightStatus, mgr *manager) {
					// Check if car still exists (it might have been deregistered by a previous green light notification)
					mgr.mu.Lock()
					var currentCar *car
					for _, c := range mgr.registeredCars {
						if c.id == id {
							currentCar = c
							break
						}
					}
					mgr.mu.Unlock()

					if currentCar == nil { // Car was already deregistered
						// fmt.Printf("Car %d callback: already deregistered.\n", id)
						return
					}

					if r.Float32() > 0.95 { // Special car logic
						fmt.Printf("Car %d: Special car, ignoring light status. Passing and deregistering.\n", id)
						mgr.DeregisterCar(id) // Special cars also leave
						return
					}

					switch status {
					case LIGHT_STATUS_RED:
						fmt.Printf("Car %d: Light is RED, I must wait.\n", id)
					case LIGHT_STATUS_GREEN:
						fmt.Printf("Car %d: Light is GREEN, I can go!\n", id)
						currentCar.run(id)    // Use currentCar which is confirmed to exist
						mgr.DeregisterCar(id) // After running on green, the car deregisters itself
					case LIGHT_STATUS_YELLOW:
						fmt.Printf("Car %d: Light is YELLOW, be careful!\n", id)
					}
				}
				m.RegisterCar(newCar) // Register the car with the manager
			}
		}
	}()

	go m.cycle(ctx) // Start the traffic light cycle

	<-ctx.Done() // Wait for interruption signal
	cancel()     // Propagate cancellation
	log.Print("Manager shutting down gracefully...")
	// Add a small delay to allow ongoing deregistration messages to print
	time.Sleep(100 * time.Millisecond)
	log.Print("Manager graceful down complete.")
}
