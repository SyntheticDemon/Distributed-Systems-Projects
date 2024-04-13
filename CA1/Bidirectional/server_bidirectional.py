# import grpc
# from concurrent import futures
# from datetime import datetime
# import bidirectional_pb2_grpc as pb2_grpc
# import bidirectional_pb2 as pb2
# import sys
# sys.path.append(r"D:\UT\Lessons\Term8\Distributed Systems\Projects")

# from CA1 import utils 


# class BidirectionalService(pb2_grpc.BidirectionalServicer):
#     def GetOrder(self, request_iterator, context):
#         available_orders = utils.read_orders_from_file("../orders.txt")
#         print('AA')
#         for request in request_iterator:
#             print('A')
#             if request.itemName in available_orders:
#                 print('B')
#                 yield pb2.ServerMessage(itemName=request.itemName,
#                                                       timestamp= str(datetime.today()))
#             else:
#                 print(f"Item not found: {request.itemName}")
        


# def serve():
#     server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
#     pb2_grpc.add_BidirectionalServicer_to_server(BidirectionalService(), server)
#     server.add_insecure_port('[::]:8081')
#     server.start()
#     server.wait_for_termination()


# if __name__ == '__main__':
#     serve()


import grpc
from concurrent import futures
from datetime import datetime
import bidirectional_pb2_grpc as pb2_grpc
import bidirectional_pb2 as pb2
import sys
sys.path.append(r"D:\UT\Lessons\Term8\Distributed Systems\Projects")

from CA1 import utils

class BidirectionalService(pb2_grpc.BidirectionalServicer):
    def GetOrder(self, request_iterator, context):
        available_orders = utils.read_orders_from_file("../orders.txt")
        response = []
        for request in request_iterator:
            for item in request.clientItems:
                if item.name in available_orders:
                    response.append(pb2.Item(name=item.name))
                else:
                    # Handle prefixed orders
                    prefixedOrders = utils.find_orders_with_prefix(str(item.name), available_orders)
                    for prefixedOrder in prefixedOrders:
                        response.append(pb2.Item(name=prefixedOrder))
            yield pb2.ServerItemList(serverItems=response, timestamp=str(datetime.today()))

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_BidirectionalServicer_to_server(BidirectionalService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()