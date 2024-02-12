let socket: WebSocket | null = null;

export function connect(username: string, token: string, callback) {
    if (!socket) {
        const host = "192.168.0.180";
        const port = 8080;
        const endpoint_base = `${host}:${port}`;
        socket = new WebSocket(`ws://${endpoint_base}/ws?username=${username}&token=${token}`);

        socket.onopen = () => {
            console.log("Successfully Connected");
        };

        socket.onmessage = msg => {
            console.log(JSON.parse(msg.data));
            callback(msg)
        };

        socket.onclose = event => {
            window.location.href = '/login';
            console.log("Socket Closed Connection: ", event);
        };

        socket.onerror = error => {
            window.location.href = '/login';
            console.log("Socket Error: ", error);
        };
    }
}

export function close() {
    if (socket) {
        socket.close()
    }
};

export function sendMsg(msg: string) {
    if (socket) {
        console.log("sending msg: ", msg);
        socket.send(msg);
    } else {
        console.error("WebSocket connection is not established.");
    }
}

export function closeWebSocketConnection() {
    if (socket) {
        socket.close();
        socket = null;
    }
}
