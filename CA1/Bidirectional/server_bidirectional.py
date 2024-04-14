import grpc
from concurrent import futures
from datetime import datetime
import bidirectional_pb2_grpc as pb2_grpc
import bidirectional_pb2 as pb2
import sys
import os

# sys.path.append(r"D:\UT\Lessons\Term8\Distributed Systems\Projects")
sys.path.append(os.path.dirname(os.path.dirname(os.getcwd())))

from CA1 import utils

class BidirectionalService(pb2_grpc.BidirectionalServicer):
    def __init__(self):
        self.available_orders = utils.read_orders_from_file("../orders.txt")

    def GetOrder(self, request_iterator, context):
        for request in request_iterator:
            response = []
            for item in request.clientItems:
                if item.name in self.available_orders:
                    response.append(pb2.Item(name=item.name))
                else:
                    prefixedOrders = utils.find_orders_with_prefix(str(item.name), self.available_orders)
                    for prefixedOrder in prefixedOrders:
                        response.append(pb2.Item(name=prefixedOrder))
            print(f"Response: {response}")
            yield pb2.ServerItemList(serverItems=response, timestamp=str(datetime.today()))
            
    

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_BidirectionalServicer_to_server(BidirectionalService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()