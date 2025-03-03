# API Gateway Documentation

## Overview

The API Gateway serves as the central entry point for handling requests from the frontend to the backend microservices in the BetterGR system. It facilitates communication between the frontend and various backend services using gRPC, ensuring efficient and scalable interactions.


The system follows a microservices architecture, with the API Gateway acting as an intermediary between the Next.js frontend and the backend microservices.
This file provides insights into the available routes, authentication mechanisms, and request/response structures.

## System Architecture

### **Diagram**:




                               +-----------------+
                               |      UI         |
                               +-----------------+
                                       |
                                       v
                            +-----------------+       +------------+
                            |   API Gateway   |------>|  Keycloak  |
                            +-----------------+       +------------+
                                      |
               ----------------------------------------------------------------
                |                 |              |             |             |
                v                 v              v             v             v
            +-----------+  +-----------+  +-----------+  +-----------+  +-----------+
            | Students  |  |  Staff    |  |  Courses  |  | Homework  |  |  Grades   |
            |    MS     |  |    MS     |  |    MS     |  |    MS     |  |    MS     |
            +-----------+  +-----------+  +-----------+  +-----------+  +-----------+
                |                 |              |             |             |
                v                 v              v             v             v
            +-----------+  +-----------+  +-----------+  +-----------+  +-----------+
            | Students  |  |  Staff    |  |  Courses  |  | Homework  |  |  Grades   |
            |    DB     |  |    DB     |  |    DB     |  |    DB     |  |    DB     |
            +-----------+  +-----------+  +-----------+  +-----------+  +-----------+


### **Key Components:**

- **Database Structure**: Each microservice has its own dedicated database.

- **API Gateway**: Central communication hub for microservices. Routes requests to the appropriate microservices using gRPC.
- **Keycloak**: Provides authentication and authorization mechanisms.

- **Microservices**:

    - Students Microservice: Manages student-related operations.

    - Courses Microservice: Handles course-related actions.

    - Homework Microservice: Manages assignments and homework submissions.

    - Grades Microservice: Handles grading operations.

    - Staff Microservice: Manages staff-related functions.



### Request Flow
Below is the detailed operations flow when the API Gateway receives a request:
1. The frontend sends HTTP requests to the API Gateway.
2. The API Gateway authenticates the request using Keycloak.

3. The API Gateway forwards the request via gRPC to the appropriate microservice.

4. The microservice processes the request and sends a response back to the API Gateway.

5. The API Gateway translates the gRPC response into an HTTP response and returns it to the frontend.

### API Endpoints

The API Gateway exposes RESTful endpoints that map to gRPC services internally. To view endpoints go to [API Endpoints](http://localhost:1234/swagger/index.html).









