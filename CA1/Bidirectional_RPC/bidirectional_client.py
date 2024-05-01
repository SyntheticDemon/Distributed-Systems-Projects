import grpc
import bidirectional_pb2
import bidirectional_pb2_grpc


def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = bidirectional_pb2_grpc.BidirectionalStub(channel)

        def request_generator(items):
            for item in items:
                yield bidirectional_pb2.ClientMessage(itemName=item)

        items = ['Potato', 'Soda', 'Sofa', 'Mango']
        responses = stub.ChatOrder(request_generator(items))
