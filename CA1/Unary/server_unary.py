import grpc
from concurrent import futures
from datetime import datetime
import unary_pb2_grpc as pb2_grpc
import unary_pb2 as pb2
import sys
import os
# sys.path.append(r"D:\UT\Lessons\Term8\Distributed Systems\Projects")
# sys.path.append("E:\\ut\\T8\\distributed\\grpc\\")
# the strucute is: sth/CA1/Unary/unary_script.py
# run this file directly from its immediate parent folder(unary) 
sys.path.append(os.path.dirname(os.path.dirname(os.getcwd())))
from CA1 import utils 




class UnaryService(pb2_grpc.UnaryServicer):
    def GetOrder(self, request, context):
        available_orders = utils.read_orders_from_file("../orders.txt")
        
        response = []
        for item in request.clientItems:
            if item.name in available_orders:
                response.append(pb2.Item(name=item.name))
            else:
                # Handel prefixed orders
                prefixedOrders = utils.find_orders_with_prefix(str(item.name), available_orders)
                for prefixedOrder in prefixedOrders:
                    response.append(pb2.Item(name=prefixedOrder))
        print(response)
        return pb2.ServerItemList(serverItems=response, timestamp= str(datetime.today()))



def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_UnaryServicer_to_server(UnaryService(), server)
    server.add_insecure_port('[::]:8081')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()