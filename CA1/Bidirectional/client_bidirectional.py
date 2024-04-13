import grpc
import bidirectional_pb2_grpc as pb2_grpc
import bidirectional_pb2 as pb2


class BidirectionalClient(object):
    """
    Client for gRPC functionality
    """
    def __init__(self):
        self.host = 'localhost'
        self.server_port = 50051

        # instantiate a channel
        self.channel = grpc.insecure_channel(
            '{}:{}'.format(self.host, self.server_port))

        # bind the client and the server
        self.stub = pb2_grpc.BidirectionalStub(self.channel)


    def generate_messages(self):
        items = [
            pb2.ClientItemList(clientItems=[pb2.Item(name='Item1'), pb2.Item(name='Rice')]),
            pb2.ClientItemList(clientItems=[pb2.Item(name='Item2')])
        ]
        for item in items:
            return item


    def get_url(self):
            responses = self.stub.GetOrder(self.generate_messages())
            for response in responses:
                print(f'{response}')
            

if __name__ == '__main__':
    client = BidirectionalClient()
    
    client.get_url()
