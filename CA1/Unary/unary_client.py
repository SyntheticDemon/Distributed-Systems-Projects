import grpc
import unary_pb2_grpc as pb2_grpc
import unary_pb2 as pb2


class UnaryClient(object):
    """
    Client for gRPC functionality
    """

    def __init__(self):
        self.host = 'localhost'
        self.server_port = 8081

        # instantiate a channel
        self.channel = grpc.insecure_channel(
            '{}:{}'.format(self.host, self.server_port))

        # bind the client and the server
        self.stub = pb2_grpc.UnaryStub(self.channel)

    def get_url(self, message):
        """
        Client function to call the rpc for GetServerResponse
        """
        message = pb2.ClientItemList(clientItems=message)
        # print(f'{message}')
        return self.stub.GetOrder(message)


if __name__ == '__main__':
    client = UnaryClient()
    items = [
            pb2.Item(name='Item'),
            pb2.Item(name='Rice')
            # pb2.Item(name='Item3')
        ]
    result = client.get_url(items)
    print(f'{result}')