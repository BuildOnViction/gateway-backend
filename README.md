# Gateway

**Project description.**
Goal: Make a good template for enterprise application with
- Microservice design
- use gokit as basework
- Use grpc as default (not http)
- Support telemetry, logging, intrument for each apis, easy for auditing
- Support cli
- BDD testing
- Support JWT
- Degging with debug/pprof/


Tool
- gokit design
- GRPC
- Watermill for CQR
- Grafana + opensensus for tracking
- ENT for schema
- Easy generating endpoints, documentations


**Commands**
Warning - need to understand gokit framework design - basic layers ie transport (http, grpc), service, middleware, endpoints
- To add new grpc endpoints
  - Define your grpc interface in api/proto/bridge
  - ```make proto```
  - Define your interface for service (gokit service)
  - ```mga generate kit endpoint path_to_service_folder```
  - Now you need to wire your endpoints (gen from service interface) with transportation layer thru middleware and transport
