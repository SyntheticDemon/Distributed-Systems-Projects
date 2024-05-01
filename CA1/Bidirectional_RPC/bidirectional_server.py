import grpc
from concurrent import futures
from datetime import datetime
import bidirectional_pb2_grpc
import bidirectional_pb2
from CA1 import utils


class BidirectionalServicer(bidirectional_pb2_grpc.BidirectionalServicer):

    def ChatOrder(self, request_iterator, context):
        available_orders = utils.read_orders_from_file("../orders.txt")

        for request in request_iterator:
            if request.itemName in available_orders:
                yield bidirectional_pb2.ServerMessage(itemName=request.itemName,
                                                      timestamp=datetime.utcnow().isoformat())
            else:
                print(f"Item not found: {request.itemName}")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    bidirectional_pb2_grpc.add_BidirectionalServicer_to_server(BidirectionalServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started at [::]:50051")
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
