import grpc
import server_streaming_pb2_grpc as pb2_grpc
import server_streaming_pb2 as pb2

class ServerStreamingClient(object):
    def __init__(self):
        self.host = 'localhost'
        self.server_port = 8081
        self.channel = grpc.insecure_channel(f'{self.host}:{self.server_port}')
        self.stub = pb2_grpc.ServerStreamingStub(self.channel)

    def get_orders(self, message):
        responses = self.stub.GetOrder(message)
        for response in responses:
            print("Response :")
            print(f'{response}')

if __name__ == '__main__':
    client = ServerStreamingClient()
    message = pb2.ClientItemList(clientItems=[pb2.Item(name='Item1'), pb2.Item(name='Rice')])
    client.get_orders(message)
    message = pb2.ClientItemList(clientItems=[pb2.Item(name='Item')])
    client.get_orders(message)