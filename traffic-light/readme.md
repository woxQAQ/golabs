# Traffic Light Control System Simulation

This is a simple command-line program that simulates a traffic light control system for a single intersection. Vehicles arrive随机ly at the intersection and react to the traffic light signals (Red, Yellow, Green).

## Features

*   **Traffic Light Cycle**: The traffic light automatically cycles through Red -> Green -> Yellow sequences. The duration for each state is configurable.
*   **Vehicle Arrival**: Vehicles arrive at the intersection at random intervals.
*   **Vehicle Behavior**:
    *   At a red light, vehicles wait.
    *   At a green light, vehicles pass through the intersection and are deregistered from the system afterward.
    *   At a yellow light, vehicles typically prepare to stop (in the current implementation, newly arriving vehicles during a yellow light will wait for the next red or green signal).
*   **Observer-like Pattern**: The system uses a mechanism similar to the Observer pattern, where the traffic light `manager` notifies registered `car` instances of light status changes.
*   **Sequential Arrival**: Vehicles are queued in the order of their arrival and are expected to pass in this order (though concurrent notification might lead to slight variations in actual execution order, notifications are based on a snapshot of an ordered list).
*   **Graceful Shutdown**: The program can be shut down gracefully via `Ctrl+C` (SIGINT) or `kill` (SIGTERM) signals.

## How to Run

1.  Ensure you have Go installed (recommended version 1.18+).
2.  Clone or download the code to your local machine.
3.  Open a terminal and navigate to the project directory (`traffic-light`).
4.  Run the following command:

    ```bash
    go run main.go
    ```

5.  You will see log messages in the terminal indicating traffic light changes and vehicle arrivals, waiting, passing, and deregistration.
6.  Press `Ctrl+C` to stop the simulation.

## Code Structure Overview

*   **`main.go`**: Contains the main logic of the program.

### Key Types and Components

*   **`lightStatus`**: Represents the state of the traffic light (`LIGHT_STATUS_RED`, `LIGHT_STATUS_GREEN`, `LIGHT_STATUS_YELLOW`).
*   **`car`**: Represents a vehicle.
    *   `id`: A unique identifier for the car.
    *   `callback`: A function invoked by the `manager` when the traffic light status changes. This callback defines how the car reacts to different light signals and is responsible for calling `manager.DeregisterCar` after passing on a green light.
    *   `run()`: Simulates the car passing through the intersection.
*   **`manager`**: The traffic light controller.
    *   `currentStatus`: The current traffic light state (using `atomic.Value` for concurrent-safe access).
    *   `redDuration`, `greenDuration`, `yellowDuration`: Durations for each light state.
    *   `mu`: A mutex to protect shared data (like `registeredCars`).
    *   `registeredCars`: A `[]*car` slice, used as a queue to store all registered vehicles waiting for or observing the traffic light. Vehicles are added to this queue in order of arrival.
    *   `nextCarID`: Used to generate unique car IDs (using `atomic.Int32`).
    *   `newManager()`: Creates and initializes a new `manager` instance.
    *   `RegisterCar(c *car)`: Registers a vehicle with the system, adding it to the end of the `registeredCars` queue, and immediately notifies the car of the current traffic light status.
    *   `DeregisterCar(carID int)`: Removes a vehicle from the `registeredCars` queue by its ID.
    *   `cycle(ctx context.Context)`: Runs the main traffic light loop, switching light states based on configured durations and notifying vehicles via `notifyCars`.
    *   `notifyCars(status lightStatus)`: Notifies all vehicles currently in the `registeredCars` queue when the light status changes. It creates a snapshot of the vehicle queue for iteration to allow cars to safely deregister themselves within their callbacks.
    *   `getStatus()`: Retrieves the current traffic light status.

### `main()` Function

*   Initializes the `manager`.
*   Sets up a goroutine to periodically and randomly generate new `car` instances and register them with the `manager`.
    *   The callback function for each newly created `car` defines how it responds to red, green, and yellow lights, and how it calls `manager.DeregisterCar` after passing on a green light.
*   Starts the `manager`'s `cycle` goroutine.
*   Listens for operating system interrupt signals to enable a graceful shutdown.

## Potential Future Improvements (Optional)

*   Implement stricter FIFO vehicle passage logic (e.g., serially processing vehicles from the head of the queue on a green light).
*   Support for multiple lanes or more complex intersections.
*   Dynamic adjustment of traffic light timings based on traffic flow.
*   Addition of unit and integration tests.
*   Splitting `car` and `manager` definitions into separate files (e.g., `types.go`) for better project structure.