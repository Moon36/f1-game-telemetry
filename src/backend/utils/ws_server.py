"""
A WebSocket server for real-time communication.
"""
import json

import websockets


def __is_valid_json__(message: str) -> bool:
    """
    Check if a message is valid JSON.

    :param message: The message to check.
    :return: True if valid JSON, False otherwise.
    """
    try:
        json.loads(message)
        return True
    except ValueError:
        return False


class WebSocketServer:
    """A WebSocket server class."""
    def __init__(self, host='localhost', port=8765):
        self.host = host
        self.port = port
        self.clients = set()
        self.server = None


    async def start(self):
        """Start the WebSocket server."""
        self.server = await websockets.serve(self.client_connect, self.host, self.port)


    async def close(self):
        """Stop the WebSocket server."""
        for client in self.clients:
            client.close()
        self.clients.clear()
        if self.server:
            self.server.close()


    async def client_connect(self, websocket: websockets.ServerConnection):
        """
        Handle new client connection.
        
        :param websocket: The new client websocket.
        """
        self.clients.add(websocket)
        print(f"Client connected: {websocket.remote_address}")

        try:
            # Keep the connection alive and handle messages
            async for message in websocket:
                print(f"Received from client: {message}")
        except websockets.exceptions.ConnectionClosed:
            print(f"Client disconnected: {websocket.remote_address}")
        finally:
            await self.client_disconnect(websocket)


    async def client_disconnect(self, websocket):
        """
        Handle client disconnection.

        :param websocket: The disconnected client websocket.
        """
        if websocket in self.clients:
            self.clients.remove(websocket)
            print(f"Client disconnected: {websocket.remote_address}")


    async def broadcast_message(self, message: str):
        """
        Broadcast a message to all connected clients.

        :param message: The message to broadcast.
        :raises ValueError: If the message is not valid JSON.
        """
        if not __is_valid_json__(message):
            raise ValueError("Invalid JSON message")

        if not self.clients:
            return

        # Create a copy of clients to avoid modification during iteration
        clients_copy = self.clients.copy()

        for client in clients_copy:
            try:
                await client.send(message)
            except websockets.exceptions.ConnectionClosed:
                # Remove disconnected clients
                await self.client_disconnect(client)
            except websockets.exceptions.WebSocketException as e:
                print(f"Error sending to client: {e}")
                await self.client_disconnect(client)
